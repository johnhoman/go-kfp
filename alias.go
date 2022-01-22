package kfp

import (
	"github.com/johnhoman/go-kfp/pipelines"
)

type (
	Pipelines     = pipelines.Interface
	GetOptions    = pipelines.GetOptions
	CreateOptions = pipelines.CreateOptions
	DeleteOptions = pipelines.DeleteOptions
	UpdateOptions = pipelines.UpdateOptions

	GetVersionOptions    = pipelines.GetOptions
	CreateVersionOptions = pipelines.CreateVersionOptions
	DeleteVersionOptions = pipelines.DeleteOptions

	Pipeline = pipelines.Pipeline
	PipelineVersion = pipelines.PipelineVersion
)

var (
	// New returns a new pipeline api client
	New = pipelines.New

	// IsNotFound returns true if the given error is not found and false otherwise
	IsNotFound = pipelines.IsNotFound

	// IsConflict returns true if an api resource with the same name already exists
	IsConflict = pipelines.IsConflict
)
