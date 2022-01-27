package pipelines

import (
    "github.com/go-openapi/runtime"
    "github.com/go-openapi/strfmt"
    "github.com/johnhoman/go-kfp/api/job/client/job_service"
    ps "github.com/johnhoman/go-kfp/api/pipeline/client/pipeline_service"
    up "github.com/johnhoman/go-kfp/api/pipeline_upload/client/pipeline_upload_service"
)

type UploadService = up.ClientService
type Service = ps.ClientService
type JobService = job_service.ClientService

type pipelineService struct {
    UploadService
    Service
    JobService
}

func NewPipelineService(transport runtime.ClientTransport) PipelineService {
    return &pipelineService{
        UploadService: up.New(transport, strfmt.Default),
        Service: ps.New(transport, strfmt.Default),
        JobService: job_service.New(transport, strfmt.Default),
    }
}