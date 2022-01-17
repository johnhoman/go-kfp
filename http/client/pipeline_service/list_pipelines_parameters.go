// Code generated by go-swagger; DO NOT EDIT.

package pipeline_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewListPipelinesParams creates a new ListPipelinesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListPipelinesParams() *ListPipelinesParams {
	return &ListPipelinesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListPipelinesParamsWithTimeout creates a new ListPipelinesParams object
// with the ability to set a timeout on a request.
func NewListPipelinesParamsWithTimeout(timeout time.Duration) *ListPipelinesParams {
	return &ListPipelinesParams{
		timeout: timeout,
	}
}

// NewListPipelinesParamsWithContext creates a new ListPipelinesParams object
// with the ability to set a context for a request.
func NewListPipelinesParamsWithContext(ctx context.Context) *ListPipelinesParams {
	return &ListPipelinesParams{
		Context: ctx,
	}
}

// NewListPipelinesParamsWithHTTPClient creates a new ListPipelinesParams object
// with the ability to set a custom HTTPClient for a request.
func NewListPipelinesParamsWithHTTPClient(client *http.Client) *ListPipelinesParams {
	return &ListPipelinesParams{
		HTTPClient: client,
	}
}

/* ListPipelinesParams contains all the parameters to send to the API endpoint
   for the list pipelines operation.

   Typically these are written to a http.Request.
*/
type ListPipelinesParams struct {

	/* Filter.

	     A url-encoded, JSON-serialized Filter protocol buffer (see
	[filter.proto](https://github.com/kubeflow/pipelines/blob/master/backend/api/filter.proto)).
	*/
	Filter *string

	/* PageSize.

	     The number of pipelines to be listed per page. If there are more pipelines
	than this number, the response message will contain a valid value in the
	nextPageToken field.

	     Format: int32
	*/
	PageSize *int32

	/* PageToken.

	     A page token to request the next page of results. The token is acquried
	from the nextPageToken field of the response from the previous
	ListPipelines call.
	*/
	PageToken *string

	/* ResourceReferenceKeyID.

	   The ID of the resource that referred to.
	*/
	ResourceReferenceKeyID *string

	/* ResourceReferenceKeyType.

	   The type of the resource that referred to.

	   Default: "UNKNOWN_RESOURCE_TYPE"
	*/
	ResourceReferenceKeyType *string

	/* SortBy.

	     Can be format of "field_name", "field_name asc" or "field_name desc"
	Ascending by default.
	*/
	SortBy *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list pipelines params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListPipelinesParams) WithDefaults() *ListPipelinesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list pipelines params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListPipelinesParams) SetDefaults() {
	var (
		resourceReferenceKeyTypeDefault = string("UNKNOWN_RESOURCE_TYPE")
	)

	val := ListPipelinesParams{
		ResourceReferenceKeyType: &resourceReferenceKeyTypeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the list pipelines params
func (o *ListPipelinesParams) WithTimeout(timeout time.Duration) *ListPipelinesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list pipelines params
func (o *ListPipelinesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list pipelines params
func (o *ListPipelinesParams) WithContext(ctx context.Context) *ListPipelinesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list pipelines params
func (o *ListPipelinesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list pipelines params
func (o *ListPipelinesParams) WithHTTPClient(client *http.Client) *ListPipelinesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list pipelines params
func (o *ListPipelinesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFilter adds the filter to the list pipelines params
func (o *ListPipelinesParams) WithFilter(filter *string) *ListPipelinesParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the list pipelines params
func (o *ListPipelinesParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithPageSize adds the pageSize to the list pipelines params
func (o *ListPipelinesParams) WithPageSize(pageSize *int32) *ListPipelinesParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the list pipelines params
func (o *ListPipelinesParams) SetPageSize(pageSize *int32) {
	o.PageSize = pageSize
}

// WithPageToken adds the pageToken to the list pipelines params
func (o *ListPipelinesParams) WithPageToken(pageToken *string) *ListPipelinesParams {
	o.SetPageToken(pageToken)
	return o
}

// SetPageToken adds the pageToken to the list pipelines params
func (o *ListPipelinesParams) SetPageToken(pageToken *string) {
	o.PageToken = pageToken
}

// WithResourceReferenceKeyID adds the resourceReferenceKeyID to the list pipelines params
func (o *ListPipelinesParams) WithResourceReferenceKeyID(resourceReferenceKeyID *string) *ListPipelinesParams {
	o.SetResourceReferenceKeyID(resourceReferenceKeyID)
	return o
}

// SetResourceReferenceKeyID adds the resourceReferenceKeyId to the list pipelines params
func (o *ListPipelinesParams) SetResourceReferenceKeyID(resourceReferenceKeyID *string) {
	o.ResourceReferenceKeyID = resourceReferenceKeyID
}

// WithResourceReferenceKeyType adds the resourceReferenceKeyType to the list pipelines params
func (o *ListPipelinesParams) WithResourceReferenceKeyType(resourceReferenceKeyType *string) *ListPipelinesParams {
	o.SetResourceReferenceKeyType(resourceReferenceKeyType)
	return o
}

// SetResourceReferenceKeyType adds the resourceReferenceKeyType to the list pipelines params
func (o *ListPipelinesParams) SetResourceReferenceKeyType(resourceReferenceKeyType *string) {
	o.ResourceReferenceKeyType = resourceReferenceKeyType
}

// WithSortBy adds the sortBy to the list pipelines params
func (o *ListPipelinesParams) WithSortBy(sortBy *string) *ListPipelinesParams {
	o.SetSortBy(sortBy)
	return o
}

// SetSortBy adds the sortBy to the list pipelines params
func (o *ListPipelinesParams) SetSortBy(sortBy *string) {
	o.SortBy = sortBy
}

// WriteToRequest writes these params to a swagger request
func (o *ListPipelinesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Filter != nil {

		// query param filter
		var qrFilter string

		if o.Filter != nil {
			qrFilter = *o.Filter
		}
		qFilter := qrFilter
		if qFilter != "" {

			if err := r.SetQueryParam("filter", qFilter); err != nil {
				return err
			}
		}
	}

	if o.PageSize != nil {

		// query param page_size
		var qrPageSize int32

		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt32(qrPageSize)
		if qPageSize != "" {

			if err := r.SetQueryParam("page_size", qPageSize); err != nil {
				return err
			}
		}
	}

	if o.PageToken != nil {

		// query param page_token
		var qrPageToken string

		if o.PageToken != nil {
			qrPageToken = *o.PageToken
		}
		qPageToken := qrPageToken
		if qPageToken != "" {

			if err := r.SetQueryParam("page_token", qPageToken); err != nil {
				return err
			}
		}
	}

	if o.ResourceReferenceKeyID != nil {

		// query param resource_reference_key.id
		var qrResourceReferenceKeyID string

		if o.ResourceReferenceKeyID != nil {
			qrResourceReferenceKeyID = *o.ResourceReferenceKeyID
		}
		qResourceReferenceKeyID := qrResourceReferenceKeyID
		if qResourceReferenceKeyID != "" {

			if err := r.SetQueryParam("resource_reference_key.id", qResourceReferenceKeyID); err != nil {
				return err
			}
		}
	}

	if o.ResourceReferenceKeyType != nil {

		// query param resource_reference_key.type
		var qrResourceReferenceKeyType string

		if o.ResourceReferenceKeyType != nil {
			qrResourceReferenceKeyType = *o.ResourceReferenceKeyType
		}
		qResourceReferenceKeyType := qrResourceReferenceKeyType
		if qResourceReferenceKeyType != "" {

			if err := r.SetQueryParam("resource_reference_key.type", qResourceReferenceKeyType); err != nil {
				return err
			}
		}
	}

	if o.SortBy != nil {

		// query param sort_by
		var qrSortBy string

		if o.SortBy != nil {
			qrSortBy = *o.SortBy
		}
		qSortBy := qrSortBy
		if qSortBy != "" {

			if err := r.SetQueryParam("sort_by", qSortBy); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
