// Code generated by go-swagger; DO NOT EDIT.

package experiment_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new experiment service API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for experiment service API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	ArchiveExperiment(params *ArchiveExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ArchiveExperimentOK, error)

	CreateExperiment(params *CreateExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateExperimentOK, error)

	DeleteExperiment(params *DeleteExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteExperimentOK, error)

	GetExperiment(params *GetExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetExperimentOK, error)

	ListExperiment(params *ListExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListExperimentOK, error)

	UnarchiveExperiment(params *UnarchiveExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*UnarchiveExperimentOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  ArchiveExperiment archives an experiment and the experiment s runs and jobs
*/
func (a *Client) ArchiveExperiment(params *ArchiveExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ArchiveExperimentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewArchiveExperimentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ArchiveExperiment",
		Method:             "POST",
		PathPattern:        "/apis/v1beta1/experiments/{id}:archive",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ArchiveExperimentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ArchiveExperimentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ArchiveExperimentDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  CreateExperiment creates a new experiment
*/
func (a *Client) CreateExperiment(params *CreateExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateExperimentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateExperimentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "CreateExperiment",
		Method:             "POST",
		PathPattern:        "/apis/v1beta1/experiments",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &CreateExperimentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateExperimentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*CreateExperimentDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  DeleteExperiment deletes an experiment without deleting the experiment s runs and jobs to avoid unexpected behaviors delete an experiment s runs and jobs before deleting the experiment
*/
func (a *Client) DeleteExperiment(params *DeleteExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteExperimentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteExperimentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteExperiment",
		Method:             "DELETE",
		PathPattern:        "/apis/v1beta1/experiments/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DeleteExperimentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteExperimentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DeleteExperimentDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetExperiment finds a specific experiment by ID
*/
func (a *Client) GetExperiment(params *GetExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetExperimentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetExperimentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetExperiment",
		Method:             "GET",
		PathPattern:        "/apis/v1beta1/experiments/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &GetExperimentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetExperimentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetExperimentDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListExperiment finds all experiments supports pagination and sorting on certain fields
*/
func (a *Client) ListExperiment(params *ListExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListExperimentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListExperimentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ListExperiment",
		Method:             "GET",
		PathPattern:        "/apis/v1beta1/experiments",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ListExperimentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListExperimentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListExperimentDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  UnarchiveExperiment restores an archived experiment the experiment s archived runs and jobs will stay archived
*/
func (a *Client) UnarchiveExperiment(params *UnarchiveExperimentParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*UnarchiveExperimentOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUnarchiveExperimentParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "UnarchiveExperiment",
		Method:             "POST",
		PathPattern:        "/apis/v1beta1/experiments/{id}:unarchive",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UnarchiveExperimentReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*UnarchiveExperimentOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*UnarchiveExperimentDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}