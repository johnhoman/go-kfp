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

// NewDisableJobParams creates a new DisableJobParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDisableJobParams() *DisableJobParams {
	return &DisableJobParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDisableJobParamsWithTimeout creates a new DisableJobParams object
// with the ability to set a timeout on a request.
func NewDisableJobParamsWithTimeout(timeout time.Duration) *DisableJobParams {
	return &DisableJobParams{
		timeout: timeout,
	}
}

// NewDisableJobParamsWithContext creates a new DisableJobParams object
// with the ability to set a context for a request.
func NewDisableJobParamsWithContext(ctx context.Context) *DisableJobParams {
	return &DisableJobParams{
		Context: ctx,
	}
}

// NewDisableJobParamsWithHTTPClient creates a new DisableJobParams object
// with the ability to set a custom HTTPClient for a request.
func NewDisableJobParamsWithHTTPClient(client *http.Client) *DisableJobParams {
	return &DisableJobParams{
		HTTPClient: client,
	}
}

/* DisableJobParams contains all the parameters to send to the API endpoint
   for the disable job operation.

   Typically these are written to a http.Request.
*/
type DisableJobParams struct {

	/* ID.

	   The ID of the job to be disabled
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the disable job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DisableJobParams) WithDefaults() *DisableJobParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the disable job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DisableJobParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the disable job params
func (o *DisableJobParams) WithTimeout(timeout time.Duration) *DisableJobParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the disable job params
func (o *DisableJobParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the disable job params
func (o *DisableJobParams) WithContext(ctx context.Context) *DisableJobParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the disable job params
func (o *DisableJobParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the disable job params
func (o *DisableJobParams) WithHTTPClient(client *http.Client) *DisableJobParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the disable job params
func (o *DisableJobParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the disable job params
func (o *DisableJobParams) WithID(id string) *DisableJobParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the disable job params
func (o *DisableJobParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DisableJobParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
