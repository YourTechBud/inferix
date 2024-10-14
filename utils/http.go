package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// HTTPResponse is a generic response for HTTP requests
type HTTPResponse[T, E any] struct {
	Status int
	Data   *T
	Error  *E
}

// MakeHTTPRequest makes an HTTP request
func MakeHTTPRequest[T, E any](ctx context.Context, method, url string, body any) (HTTPResponse[T, E], error) {
	// Marshal the body if provided
	var bodyReader io.Reader
	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return HTTPResponse[T, E]{}, err
		}
		bodyReader = bytes.NewReader(bodyData)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return HTTPResponse[T, E]{}, err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return HTTPResponse[T, E]{}, err
	}

	defer res.Body.Close()

	// Check if the status code is not 200
	if res.StatusCode >= 400 || res.StatusCode < 200 {
		data := new(E)
		if err := json.NewDecoder(res.Body).Decode(data); err != nil {
			return HTTPResponse[T, E]{}, err
		}

		return HTTPResponse[T, E]{
			Status: res.StatusCode,
			Error:  data,
		}, nil
	}

	// Decode the response
	data := new(T)
	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		return HTTPResponse[T, E]{}, err
	}

	return HTTPResponse[T, E]{
		Status: res.StatusCode,
		Data:   data,
	}, nil
}

// HasError checks if the response has an error
func (r HTTPResponse[T, E]) HasError() bool {
	return r.Error != nil
}

func WriteJSONError(w http.ResponseWriter, err error) {
	// First set the headers
	w.Header().Set("Content-Type", "application/json")

	// Check if the error is a standard error
	var standardError StandardError
	if !errors.As(err, &standardError) {
		standardError = NewStandardError(http.StatusInternalServerError, err.Error(), "internal-error")
	}

	// Set the status code
	w.WriteHeader(standardError.Status)

	// Write the error
	_ = json.NewEncoder(w).Encode(standardError)
	return
}

func WriteJSON(w http.ResponseWriter, data any) {
	// First set the headers
	w.Header().Set("Content-Type", "application/json")

	// Write the data
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
	return
}
