// Package models is the package that contains
// all the structs that will be returned in JSON.
package jsonmodels

// ErrorResponse is the struct used
// to return an error message after an error
// occured during an API call.
type ErrorResponse struct {
	Message string `json:"message"`
}
