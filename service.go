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
		Service:       ps.New(transport, strfmt.Default),
		JobService:    job_service.New(transport, strfmt.Default),
	}
}
