/*
Copyright 2022 John Homan.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kfp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/johnhoman/go-kfp/api/pipeline_upload/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-openapi/runtime"

	"github.com/johnhoman/go-kfp/api/job/client/job_service"
	jobmodels "github.com/johnhoman/go-kfp/api/job/models"
	"github.com/johnhoman/go-kfp/api/pipeline/client/pipeline_service"
	pipelinemodels "github.com/johnhoman/go-kfp/api/pipeline/models"
	"github.com/johnhoman/go-kfp/api/pipeline_upload/client/pipeline_upload_service"
)

type pipelinesApi struct {
	service  PipelineService
	authInfo runtime.ClientAuthInfoWriter
}

// CreateJob starts a scheduled job on Kubeflow. Start and end dates can be specified
// but if not provided the job will start at the current time and an end time won't
// be sent to Kubeflow so that the Kubeflow api server can pick a sensible default.
func (p *pipelinesApi) CreateJob(ctx context.Context, options *CreateJobOptions) (*Job, error) {
	start := time.Now()
	if options.StartTime != nil {
		start = *options.StartTime
	}

	body := &jobmodels.APIJob{
		Name:        options.Name,
		Description: options.Description,
		PipelineSpec: &jobmodels.APIPipelineSpec{
			Parameters: nil,
			PipelineID: options.PipelineID,
		},
		ResourceReferences: []*jobmodels.APIResourceReference{{
			Key: &jobmodels.APIResourceKey{
				ID:   options.VersionID,
				Type: jobmodels.NewAPIResourceType(jobmodels.APIResourceTypePIPELINEVERSION),
			},
			Relationship: jobmodels.NewAPIRelationship(jobmodels.APIRelationshipCREATOR),
		}},
		Trigger: &jobmodels.APITrigger{
			CronSchedule: &jobmodels.APICronSchedule{
				Cron:      options.CronSchedule,
				StartTime: strfmt.DateTime(start),
			},
		},
		NoCatchup:      true,
		MaxConcurrency: strconv.Itoa(options.MaxConcurrency),
		Enabled:        options.Enabled,
	}
	if options.EndTime != nil {
		body.Trigger.CronSchedule.EndTime = strfmt.DateTime(*options.EndTime)
	}
	out, err := p.service.CreateJob(&job_service.CreateJobParams{Body: body, Context: ctx}, p.authInfo)
	if err != nil {
		return &Job{}, err
	}
	return p.GetJob(ctx, &GetOptions{ID: out.GetPayload().ID})
}

func (p *pipelinesApi) GetJob(ctx context.Context, options *GetOptions) (*Job, error) {

	if len(options.ID) == 0 {
		// Get the ID
		if len(options.Name) == 0 {
			return &Job{}, fmt.Errorf("must specify either name or ID")
		}

		filter := map[string]interface{}{
			"predicates": []interface{}{
				map[string]interface{}{
					"key":          "name",
					"op":           "EQUALS",
					"string_value": options.Name,
				},
			},
		}
		raw, err := json.Marshal(filter)
		if err != nil {
			return &Job{}, err
		}
		in := &job_service.ListJobsParams{
			Context:  ctx,
			Filter:   stringPointer(string(raw)),
			PageSize: int32Pointer(1),
		}
		out, err := p.service.ListJobs(in, p.authInfo)
		if err != nil {
			e, ok := err.(*job_service.ListJobsDefault)
			if ok && e.Code() == http.StatusNotFound {
				return &Job{}, NewNotFound()
			}
			return &Job{}, err
		}
		if out.GetPayload().TotalSize < 1 {
			return &Job{}, NewNotFound()
		}
		job := out.GetPayload().Jobs[0]
		options = &GetOptions{ID: job.ID}
	}

	in := &job_service.GetJobParams{ID: options.ID, Context: ctx}
	out, err := p.service.GetJob(in, p.authInfo)
	if err != nil {
		e, ok := err.(*job_service.GetJobDefault)
		if ok && e.Code() == http.StatusNotFound {
			return &Job{}, NewNotFound()
		}
		return &Job{}, err
	}
	job := &Job{}
	job.Enabled = out.GetPayload().Enabled
	job.CronSchedule = out.GetPayload().Trigger.CronSchedule.Cron
	job.MaxConcurrency, err = strconv.Atoi(out.GetPayload().MaxConcurrency)
	if err != nil {
		return &Job{}, err
	}
	job.ID = out.GetPayload().ID
	job.Name = out.GetPayload().Name
	job.PipelineID = out.GetPayload().PipelineSpec.PipelineID
	job.CreatedAt = time.Time(out.GetPayload().CreatedAt)
	if len(out.GetPayload().ResourceReferences) > 0 {
		for _, ref := range out.GetPayload().ResourceReferences {
			if *ref.Key.Type == jobmodels.APIResourceTypeEXPERIMENT {
				job.ExperimentID = ref.Key.ID
			}
			if *ref.Key.Type == jobmodels.APIResourceTypePIPELINEVERSION {
				job.VersionID = ref.Key.ID
			}
		}
	}
	job.Description = out.GetPayload().Description
	job.StartTime = time.Time(out.GetPayload().Trigger.CronSchedule.StartTime)
	job.EndTime = time.Time(out.GetPayload().Trigger.CronSchedule.EndTime)

	if len(job.VersionID) > 0 {
		version, err := p.GetVersion(ctx, &GetVersionOptions{ID: job.VersionID})
		if err != nil {
			return &Job{}, err
		}
		if version.Parameters != nil {
			job.Parameters = version.Parameters
		}
	}
	return job, nil
}

// DeleteJob removes a pipeline recurring job from Kubeflow
func (p *pipelinesApi) DeleteJob(ctx context.Context, options *DeleteOptions) error {
	_, err := p.service.DeleteJob(&job_service.DeleteJobParams{ID: options.ID, Context: ctx}, p.authInfo)
	if err != nil {
		return err
	}
	return nil
}

func (p *pipelinesApi) getPipelineVersionByName(ctx context.Context, name string, pipelineId string) (*pipelinemodels.APIPipelineVersion, error) {
	rv := &pipelinemodels.APIPipelineVersion{}
	predicates := map[string]interface{}{
		"predicates": []interface{}{
			map[string]interface{}{
				"op":           "EQUALS",
				"key":          "name",
				"string_value": name,
			},
		},
	}
	raw, err := json.Marshal(predicates)
	if err != nil {
		return rv, err
	}

	// Make sure pipeline version name is unique
	versions, err := p.service.ListPipelineVersions(&pipeline_service.ListPipelineVersionsParams{
		Filter:          stringPointer(string(raw)),
		PageSize:        int32Pointer(1),
		ResourceKeyType: stringPointer(string(models.APIResourceTypePIPELINE)),
		ResourceKeyID:   stringPointer(pipelineId),
		Context:         ctx,
	}, p.authInfo)
	if err != nil {
		e, ok := err.(*pipeline_service.ListPipelineVersionsDefault)
		if ok && e.Code() == http.StatusNotFound {
			return rv, NewNotFound()
		}
		return rv, err
	}
	if len(versions.GetPayload().Versions) == 1 {
		*rv = *versions.GetPayload().Versions[0]
		return rv, nil
	}
	return rv, NewNotFound()
}

// CreateVersion creates a new version of the specified pipeline using
// the provided workflow spec.
func (p *pipelinesApi) CreateVersion(ctx context.Context, options *CreateVersionOptions) (*PipelineVersion, error) {
	rv := &PipelineVersion{}

	// Make sure pipeline exists
	if _, err := p.Get(ctx, &GetOptions{ID: options.PipelineID}); err != nil {
		return rv, err
	}

	_, err := p.getPipelineVersionByName(ctx, options.Name, options.PipelineID)
	if err != nil && !IsNotFound(err) {
		return rv, err
	}
	if err == nil {
		return rv, NewConflict()
	}

	raw, err := json.Marshal(options.Workflow)
	if err != nil {
		return rv, err
	}

	reader := runtime.NamedReader(options.Name+".yaml", bytes.NewReader(raw))
	defer func() {
		if err := reader.Close(); err != nil {
			panic("do i need to close this?" + err.Error())
		}
	}()

	version, err := p.service.UploadPipelineVersion(&pipeline_upload_service.UploadPipelineVersionParams{
		Description: stringPointer(options.Description),
		Name:        stringPointer(options.Name),
		Pipelineid:  stringPointer(options.PipelineID),
		Uploadfile:  reader,
		Context:     ctx,
	}, p.authInfo)
	if err != nil {
		return rv, err
	}
	return p.GetVersion(ctx, &GetVersionOptions{ID: version.GetPayload().ID})
}

func (p *pipelinesApi) DeleteVersion(ctx context.Context, options *DeleteOptions) error {
	_, err := p.GetVersion(ctx, &GetVersionOptions{ID: options.ID})
	if err != nil {
		return err
	}
	_, err = p.service.DeletePipelineVersion(&pipeline_service.DeletePipelineVersionParams{
		VersionID: options.ID,
		Context:   ctx,
	}, p.authInfo)
	if err != nil {
		// Maybe wrap this
		return err
	}
	return nil
}

func (p *pipelinesApi) GetVersion(ctx context.Context, options *GetVersionOptions) (*PipelineVersion, error) {
	rv := &PipelineVersion{}

	if len(options.ID) == 0 {
		if len(options.PipelineID) == 0 || len(options.Name) == 0 {
			return rv, fmt.Errorf("must specify either version Id or both of version name and pipeline ID")
		}
		out, err := p.getPipelineVersionByName(ctx, options.Name, options.PipelineID)
		if err != nil {
			return rv, err
		}
		options = &GetVersionOptions{ID: out.ID}
	}

	out, err := p.service.GetPipelineVersion(&pipeline_service.GetPipelineVersionParams{
		VersionID: options.ID,
		Context:   ctx,
	}, p.authInfo)
	if err != nil {
		e, ok := err.(*pipeline_service.GetPipelineVersionDefault)
		if ok {
			if e.Code() == http.StatusNotFound {
				return rv, NewNotFound()
			}
		}
		return rv, err
	}
	version := &PipelineVersion{
		ID:        out.GetPayload().ID,
		Name:      out.GetPayload().Name,
		CreatedAt: time.Time(out.GetPayload().CreatedAt),
	}
	if len(out.GetPayload().ResourceReferences) > 0 {
		refs := out.GetPayload().ResourceReferences
		version.PipelineID = refs[0].Key.ID
	}
	if out.GetPayload().Parameters != nil && len(out.GetPayload().Parameters) > 0 {
		params := make([]Parameter, 0, len(out.GetPayload().Parameters))
		for _, param := range out.GetPayload().Parameters {
			params = append(params, Parameter{Name: param.Name, Value: param.Value})
		}
		version.Parameters = params
	}
	return version, nil
}

func (p *pipelinesApi) getPipelineByName(ctx context.Context, name string) (*pipelinemodels.APIPipeline, error) {
	predicates := map[string]interface{}{
		"predicates": []interface{}{
			map[string]interface{}{
				"op":           "EQUALS",
				"key":          "name",
				"string_value": name,
			},
		},
	}

	raw, err := json.Marshal(predicates)
	if err != nil {
		return &pipelinemodels.APIPipeline{}, err
	}

	// How do I get the ID other than listing?
	listOut, err := p.service.ListPipelines(&pipeline_service.ListPipelinesParams{
		Context:  ctx,
		PageSize: int32Pointer(1),
		Filter:   stringPointer(string(raw)),
	}, p.authInfo)
	if err != nil {
		e, ok := err.(*pipeline_service.ListPipelineVersionsDefault)
		if ok && e.Code() == http.StatusNotFound {
			return &pipelinemodels.APIPipeline{}, NewNotFound()
		}
		return &pipelinemodels.APIPipeline{}, err
	}
	if listOut.GetPayload().TotalSize < 1 {
		return &pipelinemodels.APIPipeline{}, NewNotFound()
	}
	return listOut.GetPayload().Pipelines[0], nil
}

func (p *pipelinesApi) Create(ctx context.Context, options *CreateOptions) (*Pipeline, error) {

	_, err := p.getPipelineByName(ctx, options.Name)
	if err == nil {
		return &Pipeline{}, NewConflict()
	}
	if !IsNotFound(err) {
		return &Pipeline{}, err
	}

	raw, err := json.Marshal(options.Workflow)
	if err != nil {
		return &Pipeline{}, err
	}

	reader := runtime.NamedReader(options.Name+".yaml", bytes.NewReader(raw))
	defer func() {
		if err := reader.Close(); err != nil {
			panic("do i need to close this?" + err.Error())
		}
	}()
	params := &pipeline_upload_service.UploadPipelineParams{
		Description: stringPointer(options.Description),
		Name:        stringPointer(options.Name),
		Uploadfile:  reader,
		Context:     ctx,
	}
	out, err := p.service.UploadPipeline(params, p.authInfo)
	if err != nil {
		return &Pipeline{}, err
	}
	return &Pipeline{
		ID:               out.GetPayload().ID,
		Name:             out.GetPayload().Name,
		Description:      out.GetPayload().Description,
		CreatedAt:        time.Time(out.GetPayload().CreatedAt),
		DefaultVersionID: out.GetPayload().ID,
	}, nil
}

func (p *pipelinesApi) Get(ctx context.Context, options *GetOptions) (*Pipeline, error) {

	pl := &pipelinemodels.APIPipeline{}
	rv := &Pipeline{}

	if options.ID != "" {
		out, err := p.service.GetPipeline(&pipeline_service.GetPipelineParams{
			Context: ctx,
			ID:      options.ID,
		}, nil)
		if err != nil {
			e, ok := err.(*pipeline_service.GetPipelineDefault)
			if ok {
				if e.Code() == http.StatusNotFound {
					return &Pipeline{}, NewNotFound()
				}
			}
			return rv, err
		}
		*pl = *out.GetPayload()
	} else {
		out, err := p.getPipelineByName(ctx, options.Name)
		if err != nil {
			return rv, err
		}
		model, err := p.service.GetPipeline(&pipeline_service.GetPipelineParams{
			Context: ctx,
			ID:      out.ID,
		}, nil)
		if err != nil {
			e, ok := err.(*pipeline_service.GetPipelineDefault)
			if ok && e.Code() == http.StatusNotFound {
				return &Pipeline{}, NewNotFound()
			}
			return rv, err
		}
		*pl = *(model.GetPayload())
	}

	pipeline := &Pipeline{}
	pipeline.ID = pl.ID
	pipeline.Name = pl.Name
	pipeline.Description = pl.Description
	pipeline.CreatedAt = time.Time(pl.CreatedAt)
	if pl.DefaultVersion != nil {
		pipeline.DefaultVersionID = pl.DefaultVersion.ID
	}
	return pipeline, nil
}

func (p *pipelinesApi) Update(ctx context.Context, options *UpdateOptions) (*Pipeline, error) {
	rv := &Pipeline{}
	if _, err := p.service.GetPipelineVersion(&pipeline_service.GetPipelineVersionParams{
		Context:   ctx,
		VersionID: options.DefaultVersionID,
	}, p.authInfo); err != nil {
		e, ok := err.(*pipeline_service.GetPipelineVersionDefault)
		if ok {
			if e.Code() == http.StatusNotFound {
				return rv, NewNotFound()
			}
		}
		return rv, err
	}

	if _, err := p.service.GetPipeline(&pipeline_service.GetPipelineParams{
		ID:      options.ID,
		Context: ctx,
	}, p.authInfo); err != nil {
		if e, ok := err.(*pipeline_service.GetPipelineDefault); ok {
			if e.Code() == http.StatusNotFound {
				return rv, NewNotFound()
			}
		}
		return rv, err
	}

	_, err := p.service.UpdatePipelineDefaultVersion(&pipeline_service.UpdatePipelineDefaultVersionParams{
		PipelineID: options.ID,
		VersionID:  options.DefaultVersionID,
		Context:    ctx,
	}, p.authInfo)
	if err != nil {
		// Should probably wrap this
		return rv, err
	}
	return p.Get(ctx, &GetOptions{ID: options.ID})
}

func (p *pipelinesApi) Delete(ctx context.Context, options *DeleteOptions) error {
	_, err := p.service.DeletePipeline(
		&pipeline_service.DeletePipelineParams{Context: ctx, ID: options.ID},
		p.authInfo,
	)
	if err != nil {
		e, ok := err.(*pipeline_service.DeletePipelineDefault)
		if ok && e.Code() == http.StatusNotFound {
			return NewNotFound()
		}
		return err
	}
	return nil
}

func New(service PipelineService, authInfo runtime.ClientAuthInfoWriter) *pipelinesApi {
	return &pipelinesApi{service: service, authInfo: authInfo}
}

type namespacedClient struct {
	client    Interface
	namespace string
}

func (n *namespacedClient) Create(ctx context.Context, options *CreateOptions) (*Pipeline, error) {
	opts := &CreateOptions{}
	*opts = *options
	opts.Name = n.namespace + "-" + options.Name
	pipeline, err := n.client.Create(ctx, opts)
	if err != nil {
		return &Pipeline{}, err
	}
	return n.Get(ctx, &GetOptions{ID: pipeline.ID})
}

func (n *namespacedClient) Get(ctx context.Context, options *GetOptions) (*Pipeline, error) {
	opts := &GetOptions{}
	*opts = *options
	if len(options.Name) > 0 {
		opts.Name = n.namespace + "-" + options.Name
	}
	pipeline, err := n.client.Get(ctx, opts)
	if err != nil {
		return pipeline, err
	}
	pipeline.Name = pipeline.Name[len(n.namespace)+1:]
	return pipeline, nil
}

func (n *namespacedClient) Update(ctx context.Context, options *UpdateOptions) (*Pipeline, error) {
	pipeline, err := n.client.Update(ctx, options)
	if err != nil {
		return &Pipeline{}, err
	}
	return n.Get(ctx, &GetOptions{ID: pipeline.ID})
}

func (n *namespacedClient) Delete(ctx context.Context, options *DeleteOptions) error {
	return n.client.Delete(ctx, options)
}

func (n *namespacedClient) GetVersion(ctx context.Context, options *GetVersionOptions) (*PipelineVersion, error) {
	opts := &GetVersionOptions{}
	*opts = *options
	if len(options.Name) > 0 {
		opts.Name = n.namespace + "-" + options.Name
	}
	version, err := n.client.GetVersion(ctx, opts)
	if err != nil {
		return &PipelineVersion{}, err
	}
	version.Name = version.Name[len(n.namespace)+1:]
	return version, nil
}

func (n *namespacedClient) CreateVersion(ctx context.Context, options *CreateVersionOptions) (*PipelineVersion, error) {
	opts := &CreateVersionOptions{}
	*opts = *options
	opts.Name = n.namespace + "-" + options.Name
	version, err := n.client.CreateVersion(ctx, opts)
	if err != nil {
		return &PipelineVersion{}, err
	}
	return n.GetVersion(ctx, &GetVersionOptions{ID: version.ID})
}

func (n *namespacedClient) DeleteVersion(ctx context.Context, options *DeleteOptions) error {
	return n.client.DeleteVersion(ctx, options)
}

func (n *namespacedClient) CreateJob(ctx context.Context, options *CreateJobOptions) (*Job, error) {
	opts := &CreateJobOptions{}
	*opts = *options
	opts.Name = n.namespace + "-" + options.Name
	job, err := n.client.CreateJob(ctx, opts)
	if err != nil {
		return &Job{}, err
	}
	return n.GetJob(ctx, &GetOptions{ID: job.ID})
}

func (n *namespacedClient) GetJob(ctx context.Context, options *GetOptions) (*Job, error) {
	opts := &GetOptions{}
	*opts = *options
	if len(options.Name) > 0 {
		opts.Name = n.namespace + "-" + options.Name
	}
	job, err := n.client.GetJob(ctx, opts)
	if err != nil {
		return &Job{}, err
	}
	if len(job.Name) == 0 {
		return &Job{}, fmt.Errorf("unknown error occurred - named not received from GetJob for %s", opts.Name)
	}
	job.Name = job.Name[len(n.namespace)+1:]
	return job, nil
}

func (n *namespacedClient) DeleteJob(ctx context.Context, options *DeleteOptions) error {
	return n.client.DeleteJob(ctx, options)
}

func NewNamespaced(client Interface, namespace string) Interface {
	return &namespacedClient{client: client, namespace: namespace}
}

var _ Interface = &pipelinesApi{}
var _ Interface = &namespacedClient{}

func stringPointer(s string) *string {
	return &s
}

func int32Pointer(i int32) *int32 {
	return &i
}
