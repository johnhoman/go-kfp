// Code generated by go-swagger; DO NOT EDIT.

package job_service

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
)

// NewEnableJobParams creates a new EnableJobParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewEnableJobParams() *EnableJobParams {
	return &EnableJobParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewEnableJobParamsWithTimeout creates a new EnableJobParams object
// with the ability to set a timeout on a request.
func NewEnableJobParamsWithTimeout(timeout time.Duration) *EnableJobParams {
	return &EnableJobParams{
		timeout: timeout,
	}
}

// NewEnableJobParamsWithContext creates a new EnableJobParams object
// with the ability to set a context for a request.
func NewEnableJobParamsWithContext(ctx context.Context) *EnableJobParams {
	return &EnableJobParams{
		Context: ctx,
	}
}

// NewEnableJobParamsWithHTTPClient creates a new EnableJobParams object
// with the ability to set a custom HTTPClient for a request.
func NewEnableJobParamsWithHTTPClient(client *http.Client) *EnableJobParams {
	return &EnableJobParams{
		HTTPClient: client,
	}
}

/* EnableJobParams contains all the parameters to send to the API endpoint
   for the enable job operation.

   Typically these are written to a http.Request.
*/
type EnableJobParams struct {

	/* ID.

	   The ID of the job to be enabled
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the enable job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *EnableJobParams) WithDefaults() *EnableJobParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the enable job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *EnableJobParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the enable job params
func (o *EnableJobParams) WithTimeout(timeout time.Duration) *EnableJobParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the enable job params
func (o *EnableJobParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the enable job params
func (o *EnableJobParams) WithContext(ctx context.Context) *EnableJobParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the enable job params
func (o *EnableJobParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the enable job params
func (o *EnableJobParams) WithHTTPClient(client *http.Client) *EnableJobParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the enable job params
func (o *EnableJobParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the enable job params
func (o *EnableJobParams) WithID(id string) *EnableJobParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the enable job params
func (o *EnableJobParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *EnableJobParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
