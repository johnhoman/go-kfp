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

// NewListPipelineVersionsParams creates a new ListPipelineVersionsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListPipelineVersionsParams() *ListPipelineVersionsParams {
	return &ListPipelineVersionsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListPipelineVersionsParamsWithTimeout creates a new ListPipelineVersionsParams object
// with the ability to set a timeout on a request.
func NewListPipelineVersionsParamsWithTimeout(timeout time.Duration) *ListPipelineVersionsParams {
	return &ListPipelineVersionsParams{
		timeout: timeout,
	}
}

// NewListPipelineVersionsParamsWithContext creates a new ListPipelineVersionsParams object
// with the ability to set a context for a request.
func NewListPipelineVersionsParamsWithContext(ctx context.Context) *ListPipelineVersionsParams {
	return &ListPipelineVersionsParams{
		Context: ctx,
	}
}

// NewListPipelineVersionsParamsWithHTTPClient creates a new ListPipelineVersionsParams object
// with the ability to set a custom HTTPClient for a request.
func NewListPipelineVersionsParamsWithHTTPClient(client *http.Client) *ListPipelineVersionsParams {
	return &ListPipelineVersionsParams{
		HTTPClient: client,
	}
}

/* ListPipelineVersionsParams contains all the parameters to send to the API endpoint
   for the list pipeline versions operation.

   Typically these are written to a http.Request.
*/
type ListPipelineVersionsParams struct {

	/* Filter.

	     A base-64 encoded, JSON-serialized Filter protocol buffer (see
	filter.proto).
	*/
	Filter *string

	/* PageSize.

	     The number of pipeline versions to be listed per page. If there are more
	pipeline versions than this number, the response message will contain a
	nextPageToken field you can use to fetch the next page.

	     Format: int32
	*/
	PageSize *int32

	/* PageToken.

	     A page token to request the next page of results. The token is acquried
	from the nextPageToken field of the response from the previous
	ListPipelineVersions call or can be omitted when fetching the first page.
	*/
	PageToken *string

	/* ResourceKeyID.

	   The ID of the resource that referred to.
	*/
	ResourceKeyID *string

	/* ResourceKeyType.

	   The type of the resource that referred to.

	   Default: "UNKNOWN_RESOURCE_TYPE"
	*/
	ResourceKeyType *string

	/* SortBy.

	     Can be format of "field_name", "field_name asc" or "field_name desc"
	Ascending by default.
	*/
	SortBy *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list pipeline versions params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListPipelineVersionsParams) WithDefaults() *ListPipelineVersionsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list pipeline versions params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListPipelineVersionsParams) SetDefaults() {
	var (
		resourceKeyTypeDefault = string("UNKNOWN_RESOURCE_TYPE")
	)

	val := ListPipelineVersionsParams{
		ResourceKeyType: &resourceKeyTypeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithTimeout(timeout time.Duration) *ListPipelineVersionsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithContext(ctx context.Context) *ListPipelineVersionsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithHTTPClient(client *http.Client) *ListPipelineVersionsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFilter adds the filter to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithFilter(filter *string) *ListPipelineVersionsParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithPageSize adds the pageSize to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithPageSize(pageSize *int32) *ListPipelineVersionsParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetPageSize(pageSize *int32) {
	o.PageSize = pageSize
}

// WithPageToken adds the pageToken to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithPageToken(pageToken *string) *ListPipelineVersionsParams {
	o.SetPageToken(pageToken)
	return o
}

// SetPageToken adds the pageToken to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetPageToken(pageToken *string) {
	o.PageToken = pageToken
}

// WithResourceKeyID adds the resourceKeyID to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithResourceKeyID(resourceKeyID *string) *ListPipelineVersionsParams {
	o.SetResourceKeyID(resourceKeyID)
	return o
}

// SetResourceKeyID adds the resourceKeyId to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetResourceKeyID(resourceKeyID *string) {
	o.ResourceKeyID = resourceKeyID
}

// WithResourceKeyType adds the resourceKeyType to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithResourceKeyType(resourceKeyType *string) *ListPipelineVersionsParams {
	o.SetResourceKeyType(resourceKeyType)
	return o
}

// SetResourceKeyType adds the resourceKeyType to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetResourceKeyType(resourceKeyType *string) {
	o.ResourceKeyType = resourceKeyType
}

// WithSortBy adds the sortBy to the list pipeline versions params
func (o *ListPipelineVersionsParams) WithSortBy(sortBy *string) *ListPipelineVersionsParams {
	o.SetSortBy(sortBy)
	return o
}

// SetSortBy adds the sortBy to the list pipeline versions params
func (o *ListPipelineVersionsParams) SetSortBy(sortBy *string) {
	o.SortBy = sortBy
}

// WriteToRequest writes these params to a swagger request
func (o *ListPipelineVersionsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	if o.ResourceKeyID != nil {

		// query param resource_key.id
		var qrResourceKeyID string

		if o.ResourceKeyID != nil {
			qrResourceKeyID = *o.ResourceKeyID
		}
		qResourceKeyID := qrResourceKeyID
		if qResourceKeyID != "" {

			if err := r.SetQueryParam("resource_key.id", qResourceKeyID); err != nil {
				return err
			}
		}
	}

	if o.ResourceKeyType != nil {

		// query param resource_key.type
		var qrResourceKeyType string

		if o.ResourceKeyType != nil {
			qrResourceKeyType = *o.ResourceKeyType
		}
		qResourceKeyType := qrResourceKeyType
		if qResourceKeyType != "" {

			if err := r.SetQueryParam("resource_key.type", qResourceKeyType); err != nil {
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
