package kfp

import (
    "github.com/go-openapi/runtime"
    "github.com/johnhoman/go-kfp/pipelines"
)

type (
	Pipelines = pipelines.Interface
    GetOptions = pipelines.GetOptions
    CreateOptions = pipelines.CreateOptions
    DeleteOptions = pipelines.DeleteOptions
    UpdateOptions = pipelines.UpdateOptions

    GetVersionOptions = pipelines.GetOptions
    CreateVersionOptions = pipelines.CreateVersionOptions
    DeleteVersionOptions = pipelines.DeleteOptions
)

func New(service pipelines.PipelineService, authInfo runtime.ClientAuthInfoWriter) pipelines.Interface {
    return pipelines.New(service, authInfo)
}