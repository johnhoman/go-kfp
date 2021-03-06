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
	"context"
	"github.com/johnhoman/go-kfp/api/experiment/client/experiment_service"

	"github.com/go-openapi/runtime"

	jobs "github.com/johnhoman/go-kfp/api/job/client/job_service"
	ps "github.com/johnhoman/go-kfp/api/pipeline/client/pipeline_service"
	up "github.com/johnhoman/go-kfp/api/pipeline_upload/client/pipeline_upload_service"
)

type PipelineService interface {
	DeletePipeline(params *ps.DeletePipelineParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.DeletePipelineOK, error)
	DeletePipelineVersion(params *ps.DeletePipelineVersionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.DeletePipelineVersionOK, error)
	GetPipeline(params *ps.GetPipelineParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.GetPipelineOK, error)
	GetPipelineVersion(params *ps.GetPipelineVersionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.GetPipelineVersionOK, error)
	UpdatePipelineDefaultVersion(params *ps.UpdatePipelineDefaultVersionParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.UpdatePipelineDefaultVersionOK, error)
	ListPipelines(params *ps.ListPipelinesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.ListPipelinesOK, error)
	ListPipelineVersions(params *ps.ListPipelineVersionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ps.ClientOption) (*ps.ListPipelineVersionsOK, error)

	UploadPipeline(params *up.UploadPipelineParams, authInfo runtime.ClientAuthInfoWriter, opts ...up.ClientOption) (*up.UploadPipelineOK, error)
	UploadPipelineVersion(params *up.UploadPipelineVersionParams, authInfo runtime.ClientAuthInfoWriter, opts ...up.ClientOption) (*up.UploadPipelineVersionOK, error)

	CreateJob(params *jobs.CreateJobParams, authInfo runtime.ClientAuthInfoWriter, opts ...jobs.ClientOption) (*jobs.CreateJobOK, error)
	DeleteJob(params *jobs.DeleteJobParams, authInfo runtime.ClientAuthInfoWriter, opts ...jobs.ClientOption) (*jobs.DeleteJobOK, error)
	GetJob(params *jobs.GetJobParams, authInfo runtime.ClientAuthInfoWriter, opts ...jobs.ClientOption) (*jobs.GetJobOK, error)
	ListJobs(params *jobs.ListJobsParams, authInfo runtime.ClientAuthInfoWriter, opts ...jobs.ClientOption) (*jobs.ListJobsOK, error)

	CreateExperiment(params *experiment_service.CreateExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...experiment_service.ClientOption) (*experiment_service.CreateExperimentOK, error)
	DeleteExperiment(params *experiment_service.DeleteExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...experiment_service.ClientOption) (*experiment_service.DeleteExperimentOK, error)
	GetExperiment(params *experiment_service.GetExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...experiment_service.ClientOption) (*experiment_service.GetExperimentOK, error)
	ListExperiment(params *experiment_service.ListExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...experiment_service.ClientOption) (*experiment_service.ListExperimentOK, error)
}

type Interface interface {
	Create(ctx context.Context, options *CreateOptions) (*Pipeline, error)
	Get(ctx context.Context, options *GetOptions) (*Pipeline, error)
	Update(ctx context.Context, options *UpdateOptions) (*Pipeline, error)
	Delete(ctx context.Context, options *DeleteOptions) error

	GetVersion(ctx context.Context, options *GetVersionOptions) (*PipelineVersion, error)
	CreateVersion(ctx context.Context, options *CreateVersionOptions) (*PipelineVersion, error)
	DeleteVersion(ctx context.Context, options *DeleteOptions) error

	CreateJob(ctx context.Context, options *CreateJobOptions) (*Job, error)
	GetJob(ctx context.Context, options *GetOptions) (*Job, error)
	DeleteJob(ctx context.Context, options *DeleteOptions) error

	CreateExperiment(ctx context.Context, options *CreateExperimentOptions) (*Experiment, error)
	GetExperiment(ctx context.Context, options *GetOptions) (*Experiment, error)
	DeleteExperiment(ctx context.Context, options *DeleteOptions) error
}
