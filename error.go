package routing

import "net/http"

type HTTPError interface {
	error
	StatusCode() int
}

type httpError struct {
	Status  int    `json:"status" xml:"status"`
	Message string `json:"message" xml:"message"`
}

func (e *httpError) Error() string {
	return e.Message
}
func (e *httpError) StatusCode() int {
	return e.Status
}

func (e *httpError) String() string {
	return e.Message
}
func NewHTTPError(status int, message ...string) HTTPError {
	if len(message) > 0 {
		return &httpError{Status: status, Message: message[0]}
	}
	return &httpError{Status: status, Message: http.StatusText(status)}
}
