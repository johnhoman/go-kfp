// Code generated by go-swagger; DO NOT EDIT.

package pipeline_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/johnhoman/go-kfp/api/pipeline/models"
)

// GetPipelineVersionReader is a Reader for the GetPipelineVersion structure.
type GetPipelineVersionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPipelineVersionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPipelineVersionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetPipelineVersionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetPipelineVersionOK creates a GetPipelineVersionOK with default headers values
func NewGetPipelineVersionOK() *GetPipelineVersionOK {
	return &GetPipelineVersionOK{}
}

/* GetPipelineVersionOK describes a response with status code 200, with default header values.

A successful response.
*/
type GetPipelineVersionOK struct {
	Payload *models.APIPipelineVersion
}

func (o *GetPipelineVersionOK) Error() string {
	return fmt.Sprintf("[GET /apis/v1beta1/pipeline_versions/{version_id}][%d] getPipelineVersionOK  %+v", 200, o.Payload)
}
func (o *GetPipelineVersionOK) GetPayload() *models.APIPipelineVersion {
	return o.Payload
}

func (o *GetPipelineVersionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIPipelineVersion)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPipelineVersionDefault creates a GetPipelineVersionDefault with default headers values
func NewGetPipelineVersionDefault(code int) *GetPipelineVersionDefault {
	return &GetPipelineVersionDefault{
		_statusCode: code,
	}
}

/* GetPipelineVersionDefault describes a response with status code -1, with default header values.

GetPipelineVersionDefault get pipeline version default
*/
type GetPipelineVersionDefault struct {
	_statusCode int

	Payload *models.APIStatus
}

// Code gets the status code for the get pipeline version default response
func (o *GetPipelineVersionDefault) Code() int {
	return o._statusCode
}

func (o *GetPipelineVersionDefault) Error() string {
	return fmt.Sprintf("[GET /apis/v1beta1/pipeline_versions/{version_id}][%d] GetPipelineVersion default  %+v", o._statusCode, o.Payload)
}
func (o *GetPipelineVersionDefault) GetPayload() *models.APIStatus {
	return o.Payload
}

func (o *GetPipelineVersionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
