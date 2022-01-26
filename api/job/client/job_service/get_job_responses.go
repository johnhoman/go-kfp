// Code generated by go-swagger; DO NOT EDIT.

package job_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/johnhoman/go-kfp/api/job/models"
)

// GetJobReader is a Reader for the GetJob structure.
type GetJobReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetJobReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetJobOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetJobDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetJobOK creates a GetJobOK with default headers values
func NewGetJobOK() *GetJobOK {
	return &GetJobOK{}
}

/* GetJobOK describes a response with status code 200, with default header values.

A successful response.
*/
type GetJobOK struct {
	Payload *models.APIJob
}

func (o *GetJobOK) Error() string {
	return fmt.Sprintf("[GET /apis/v1beta1/jobs/{id}][%d] getJobOK  %+v", 200, o.Payload)
}
func (o *GetJobOK) GetPayload() *models.APIJob {
	return o.Payload
}

func (o *GetJobOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIJob)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetJobDefault creates a GetJobDefault with default headers values
func NewGetJobDefault(code int) *GetJobDefault {
	return &GetJobDefault{
		_statusCode: code,
	}
}

/* GetJobDefault describes a response with status code -1, with default header values.

GetJobDefault get job default
*/
type GetJobDefault struct {
	_statusCode int

	Payload *models.APIStatus
}

// Code gets the status code for the get job default response
func (o *GetJobDefault) Code() int {
	return o._statusCode
}

func (o *GetJobDefault) Error() string {
	return fmt.Sprintf("[GET /apis/v1beta1/jobs/{id}][%d] GetJob default  %+v", o._statusCode, o.Payload)
}
func (o *GetJobDefault) GetPayload() *models.APIStatus {
	return o.Payload
}

func (o *GetJobDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
