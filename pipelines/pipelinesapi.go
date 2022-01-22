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

package pipelines

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"

	ps "github.com/johnhoman/go-kfp/api/pipeline/client/pipeline_service"
	"github.com/johnhoman/go-kfp/api/pipeline/models"
	up "github.com/johnhoman/go-kfp/api/pipeline_upload/client/pipeline_upload_service"
)

type pipelinesApi struct {
	service  PipelineService
	authInfo runtime.ClientAuthInfoWriter
}

func (p *pipelinesApi) CreateVersion(ctx context.Context, options *CreateVersionOptions) (*PipelineVersion, error) {
	rv := &PipelineVersion{}

	// Make sure pipeline exists
	if _, err := p.Get(ctx, &GetOptions{ID: options.PipelineID}); err != nil {
		return rv, err
	}

	predicates := map[string]interface{}{
		"predicates": []interface{}{
			map[string]interface{}{
				"op":           "EQUALS",
				"key":          "name",
				"string_value": options.Name,
			},
		},
	}
	raw, err := json.Marshal(predicates)
	if err != nil {
		return rv, err
	}

	// Make sure pipeline version name is unique
	versions, err := p.service.ListPipelineVersions(&ps.ListPipelineVersionsParams{
		Filter:          stringPointer(string(raw)),
		PageSize:        int32Pointer(1),
		ResourceKeyType: stringPointer(string(models.APIResourceTypePIPELINE)),
		ResourceKeyID:   stringPointer(options.PipelineID),
		Context:         ctx,
	}, p.authInfo)
	if err != nil {
		return rv, err
	}
	if len(versions.GetPayload().Versions) == 1 {
		return rv, NewConflict()
	}

	raw, err = json.Marshal(options.Workflow)
	if err != nil {
		return rv, err
	}

	reader := runtime.NamedReader(options.Name+".json", bytes.NewReader(raw))
	defer func() {
		if err := reader.Close(); err != nil {
			panic("do i need to close this?" + err.Error())
		}
	}()

	version, err := p.service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
		Description: stringPointer(options.Description),
		Name:        stringPointer(options.Name),
		Pipelineid:  stringPointer(options.PipelineID),
		Uploadfile:  reader,
		Context:     ctx,
	}, p.authInfo)
	if err != nil {
		return rv, err
	}
	return p.GetVersion(ctx, &GetOptions{ID: version.GetPayload().ID})
}

func (p *pipelinesApi) DeleteVersion(ctx context.Context, options *DeleteOptions) error {
	_, err := p.GetVersion(ctx, &GetOptions{ID: options.ID})
	if err != nil {
		return err
	}
	_, err = p.service.DeletePipelineVersion(&ps.DeletePipelineVersionParams{
		VersionID: options.ID,
		Context:   ctx,
	}, p.authInfo)
	if err != nil {
		// Maybe wrap this
		return err
	}
	return nil
}

func (p *pipelinesApi) GetVersion(ctx context.Context, options *GetOptions) (*PipelineVersion, error) {
	out, err := p.service.GetPipelineVersion(&ps.GetPipelineVersionParams{
		VersionID: options.ID,
		Context:   ctx,
	}, p.authInfo)
	if err != nil {
		e, ok := err.(*ps.GetPipelineVersionDefault)
		if ok {
			if e.Code() == http.StatusNotFound {
				return &PipelineVersion{}, NewNotFound()
			}
		}
	}
	return &PipelineVersion{
		ID:         out.GetPayload().ID,
		Name:       out.GetPayload().Name,
		CreatedAt:  time.Time(out.GetPayload().CreatedAt),
		PipelineID: options.ID,
	}, nil
}

func (p *pipelinesApi) getPipelineByName(ctx context.Context, name string) (*models.APIPipeline, error) {
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
		return &models.APIPipeline{}, err
	}

	// How do I get the ID other than listing?
	listOut, err := p.service.ListPipelines(&ps.ListPipelinesParams{
		Context:  ctx,
		PageSize: int32Pointer(1),
		Filter:   stringPointer(string(raw)),
	}, p.authInfo)
	if err != nil {
		return &models.APIPipeline{}, err
	}
	if listOut.GetPayload().TotalSize < 1 {
		return &models.APIPipeline{}, NewNotFound()
	}
	return listOut.GetPayload().Pipelines[0], nil
}

func (p *pipelinesApi) Create(ctx context.Context, options *CreateOptions) (*Pipeline, error) {

	_, err := p.getPipelineByName(ctx, options.Name)
	if err == nil {
		return &Pipeline{}, NewConflict()
	}

	raw, err := json.Marshal(options.Workflow)
	if err != nil {
		return &Pipeline{}, err
	}

	reader := runtime.NamedReader(options.Name+".json", bytes.NewReader(raw))
	defer func() {
		if err := reader.Close(); err != nil {
			panic("do i need to close this?" + err.Error())
		}
	}()
	params := &up.UploadPipelineParams{
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

	pl := &models.APIPipeline{}
	rv := &Pipeline{}

	if options.ID != "" {
		out, err := p.service.GetPipeline(&ps.GetPipelineParams{
			Context: ctx,
			ID:      options.ID,
		}, nil)
		if err != nil {
			e, ok := err.(*ps.GetPipelineDefault)
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
		model, err := p.service.GetPipeline(&ps.GetPipelineParams{
			Context: ctx,
			ID:      out.ID,
		}, nil)
		if err != nil {
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
	if _, err := p.service.GetPipelineVersion(&ps.GetPipelineVersionParams{
		Context:   ctx,
		VersionID: options.DefaultVersionID,
	}, p.authInfo); err != nil {
		e, ok := err.(*ps.GetPipelineVersionDefault)
		if ok {
			if e.Code() == http.StatusNotFound {
				return rv, NewNotFound()
			}
		}
		return rv, err
	}

	if _, err := p.service.GetPipeline(&ps.GetPipelineParams{
		ID:      options.ID,
		Context: ctx,
	}, p.authInfo); err != nil {
		if e, ok := err.(*ps.GetPipelineDefault); ok {
			if e.Code() == http.StatusNotFound {
				return rv, NewNotFound()
			}
		}
		return rv, err
	}

	_, err := p.service.UpdatePipelineDefaultVersion(&ps.UpdatePipelineDefaultVersionParams{
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
		&ps.DeletePipelineParams{Context: ctx, ID: options.ID},
		p.authInfo,
	)
	if def, ok := err.(*ps.DeletePipelineDefault); ok {
		if def.Code() == http.StatusNotFound {
			return NewNotFound()
		}
	}
	return err
}

func New(service PipelineService, authInfo runtime.ClientAuthInfoWriter) *pipelinesApi {
	return &pipelinesApi{service: service, authInfo: authInfo}
}

var _ Interface = &pipelinesApi{}

func stringPointer(s string) *string {
	return &s
}

func int32Pointer(i int32) *int32 {
	return &i
}
