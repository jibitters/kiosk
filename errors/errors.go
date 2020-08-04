package errors

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Type encapsulates a general error type that can be used in all layers.
type Type struct {
	FingerPrint    string  `json:"fingerprint"`
	Errors         []Error `json:"errors"`
	HTTPStatusCode int     `json:"status"`
}

// Error encapsulates an specific error. An error type may include two or more errors.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// String representation of Type.
func (t *Type) String() string {
	return fmt.Sprintf("%v", *t)
}

// Error implementation of error.
func (t *Type) Error() string {
	return t.String()
}

// InvalidRequestBody is a helper method that indicates the request body is not valid.
func InvalidRequestBody() *Type {
	return &Type{uuid.New().String(), []Error{{"invalid.json.format", ""}},
		http.StatusBadRequest}
}

// InvalidArgument is a helper method that indicates the provided argument is not valid.
func InvalidArgument(code, message string) *Type {
	return &Type{uuid.New().String(), []Error{{code, message}},
		http.StatusBadRequest}
}

// Unauthorized is a helper method that indicates the request is not authenticated.
func Unauthorized(message string) *Type {
	return &Type{uuid.New().String(), []Error{{"unauthorized", message}},
		http.StatusUnauthorized}
}

// NotFound is a helper method that indicates the resource not found.
func NotFound(code, message string) *Type {
	return &Type{uuid.New().String(), []Error{{code, message}},
		http.StatusNotFound}
}

// AlreadyExists is a helper method that indicates the resource already exists.
func AlreadyExists(code, message string) *Type {
	return &Type{uuid.New().String(), []Error{{code, message}},
		http.StatusPreconditionFailed}
}

// PreconditionFailed is a helper method that indicates some precondition failure.
func PreconditionFailed(code, message string) *Type {
	return &Type{uuid.New().String(), []Error{{code, message}},
		http.StatusPreconditionFailed}
}

// RequestTimeout is a helper method that indicates request timeout occurred.
func RequestTimeout(message string) *Type {
	return &Type{uuid.New().String(), []Error{{"request.timeout", message}},
		http.StatusRequestTimeout}
}

// ServiceUnavailable is a helper method that indicates the server is not available for now.
func ServiceUnavailable(message string) *Type {
	return &Type{uuid.New().String(), []Error{{"service.not_available", message}},
		http.StatusServiceUnavailable}
}

// InternalServerError is a helper method that indicates an internal server error occurred.
func InternalServerError(code, message string) *Type {
	return &Type{uuid.New().String(), []Error{{code, message}},
		http.StatusInternalServerError}
}

// NotImplemented is a helper method that indicates the service is not implemented yet.
func NotImplemented() *Type {
	return &Type{uuid.New().String(), []Error{{"service.not_implemented", ""}},
		http.StatusNotImplemented}
}
