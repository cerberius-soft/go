// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"cerberius.com/go-client/generated/models"
)

// EmailValidationRequestDataReader is a Reader for the EmailValidationRequestData structure.
type EmailValidationRequestDataReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *EmailValidationRequestDataReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewEmailValidationRequestDataOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewEmailValidationRequestDataDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewEmailValidationRequestDataOK creates a EmailValidationRequestDataOK with default headers values
func NewEmailValidationRequestDataOK() *EmailValidationRequestDataOK {
	return &EmailValidationRequestDataOK{}
}

/*
EmailValidationRequestDataOK describes a response with status code 200, with default header values.

Response for email validation
*/
type EmailValidationRequestDataOK struct {
	Payload *models.EmailLookupResponse
}

// IsSuccess returns true when this email validation request data o k response has a 2xx status code
func (o *EmailValidationRequestDataOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this email validation request data o k response has a 3xx status code
func (o *EmailValidationRequestDataOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this email validation request data o k response has a 4xx status code
func (o *EmailValidationRequestDataOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this email validation request data o k response has a 5xx status code
func (o *EmailValidationRequestDataOK) IsServerError() bool {
	return false
}

// IsCode returns true when this email validation request data o k response a status code equal to that given
func (o *EmailValidationRequestDataOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the email validation request data o k response
func (o *EmailValidationRequestDataOK) Code() int {
	return 200
}

func (o *EmailValidationRequestDataOK) Error() string {
	return fmt.Sprintf("[POST /email-lookup][%d] emailValidationRequestDataOK  %+v", 200, o.Payload)
}

func (o *EmailValidationRequestDataOK) String() string {
	return fmt.Sprintf("[POST /email-lookup][%d] emailValidationRequestDataOK  %+v", 200, o.Payload)
}

func (o *EmailValidationRequestDataOK) GetPayload() *models.EmailLookupResponse {
	return o.Payload
}

func (o *EmailValidationRequestDataOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EmailLookupResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewEmailValidationRequestDataDefault creates a EmailValidationRequestDataDefault with default headers values
func NewEmailValidationRequestDataDefault(code int) *EmailValidationRequestDataDefault {
	return &EmailValidationRequestDataDefault{
		_statusCode: code,
	}
}

/*
EmailValidationRequestDataDefault describes a response with status code -1, with default header values.

Generic error response
*/
type EmailValidationRequestDataDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this email validation request data default response has a 2xx status code
func (o *EmailValidationRequestDataDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this email validation request data default response has a 3xx status code
func (o *EmailValidationRequestDataDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this email validation request data default response has a 4xx status code
func (o *EmailValidationRequestDataDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this email validation request data default response has a 5xx status code
func (o *EmailValidationRequestDataDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this email validation request data default response a status code equal to that given
func (o *EmailValidationRequestDataDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the email validation request data default response
func (o *EmailValidationRequestDataDefault) Code() int {
	return o._statusCode
}

func (o *EmailValidationRequestDataDefault) Error() string {
	return fmt.Sprintf("[POST /email-lookup][%d] emailValidationRequestData default  %+v", o._statusCode, o.Payload)
}

func (o *EmailValidationRequestDataDefault) String() string {
	return fmt.Sprintf("[POST /email-lookup][%d] emailValidationRequestData default  %+v", o._statusCode, o.Payload)
}

func (o *EmailValidationRequestDataDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *EmailValidationRequestDataDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
