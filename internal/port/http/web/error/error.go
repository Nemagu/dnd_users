package weberror

import "fmt"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"field: %s\nmessage: %s", e.Field, e.Message,
	)
}

type ResponseError struct {
	StatusCode int               `json:"-"`
	Detail     string            `json:"detail,omitempty"`
	Errors     []ValidationError `json:"errors,omitempty"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf(
		"status code: %d\ndetail: %s\nerrors: %s", e.StatusCode, e.Detail, e.Errors,
	)
}
