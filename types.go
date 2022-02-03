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
	"time"
)

// GetOptions is data required by the pipelines api to identify a pipeline or version
// Either Name or ID can be used is applicable. If ID and Name are both provided then only
// ID will be used
type GetOptions struct {
	ID   string
	Name string
}

type GetVersionOptions struct {
	ID         string
	PipelineID string
	Name       string
}

type CreateOptions struct {
	Description string
	Name        string
	Workflow    map[string]interface{}
}

type CreateVersionOptions struct {
	Description string
	Name        string
	Workflow    map[string]interface{}
	PipelineID  string
}

type UpdateOptions struct {
	ID               string
	DefaultVersionID string
}

type DeleteOptions struct {
	ID string
}

type Pipeline struct {
	ID               string
	Name             string
	Description      string
	CreatedAt        time.Time
	DefaultVersionID string
}

type Parameter struct {
	Name  string
	Value string
}

type PipelineVersion struct {
	ID         string
	Name       string
	CreatedAt  time.Time
	PipelineID string
	Parameters []Parameter
}

type CreateJobOptions struct {
	Name           string
	Description    string
	PipelineID     string
	VersionID      string
	ExperimentID   string
	CronSchedule   string
	StartTime      *time.Time
	EndTime        *time.Time
	MaxConcurrency int
	Enabled        bool
}

type GetJobOption struct {
	Name string
	ID   string
}

type Job struct {
	ID             string
	Name           string
	Description    string
	PipelineID     string
	VersionID      string
	ExperimentID   string
	CronSchedule   string
	StartTime      time.Time
	EndTime        time.Time
	MaxConcurrency int
	Enabled        bool
	CreatedAt      time.Time
}
