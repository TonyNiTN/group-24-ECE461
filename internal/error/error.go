package error

import "fmt"

// Create base error class for sending requests
type RequestError struct {
	// HTTP status code
	Statuscode int

	// Error message
	Message string
}

// Create Error string
func (re *RequestError) Error() string {
	return re.Message
}

// Bad Request Error
func BadRequestError(requestType string, errorMessage string) *RequestError {
	return &RequestError{
		Statuscode: 400,
		Message:    fmt.Sprintf("%s Error: Bad request: %s", requestType, errorMessage),
	}
}

// Not Found Error
func NotFoundError(requestType string, query string) *RequestError {
	return &RequestError{
		Statuscode: 404,
		Message:    fmt.Sprintf("%s Error: could not find: %s", requestType, query),
	}
}



