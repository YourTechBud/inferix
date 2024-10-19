package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// HTTPResponse is a generic response for HTTP requests
type HTTPResponse[T, E any] struct {
	Status int
	Data   *T
	Error  *E
}

// HasError checks if the response has an error
func (r HTTPResponse[T, E]) HasError() bool {
	return r.Error != nil
}

// HTTPStreamingResponse is a generic response for HTTP streaming requests
type HTTPStreamingResponse struct {
	Status  int
	Message string
	Error   error
}

// MakeHTTPRequest makes an HTTP request
func MakeHTTPRequest[T, E any](ctx context.Context, method, url string, body any) (HTTPResponse[T, E], error) {
	// Prepare the request
	req, err := prepareHTTPRequest(ctx, method, url, body)
	if err != nil {
		return HTTPResponse[T, E]{}, err
	}

	// Fire the request
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

func MakeHTTPStream(ctx context.Context, method, url string, body any) func(yield func(element HTTPStreamingResponse) bool) {
	return func(yield func(element HTTPStreamingResponse) bool) {
		// Prepare the request
		req, err := prepareHTTPRequest(ctx, method, url, body)
		if err != nil {
			yield(HTTPStreamingResponse{Error: err})
			return
		}

		// Fire the request
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			yield(HTTPStreamingResponse{Error: err})
			return
		}

		defer res.Body.Close()

		// Check if the status code is not 2xx
		if res.StatusCode >= 400 || res.StatusCode < 200 {
			data, _ := io.ReadAll(res.Body)
			yield(HTTPStreamingResponse{Message: string(data), Status: res.StatusCode})
			return
		}

		// Decode the response
		reader := bufio.NewReader(res.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				yield(HTTPStreamingResponse{Error: err})
				return
			}

			// Skip empty lines
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if !yield(HTTPStreamingResponse{Message: line, Status: res.StatusCode}) {
				return
			}
		}
	}
}

func prepareHTTPRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	// Marshal the body if provided
	var bodyReader io.Reader
	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyData)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	return req, nil
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
}

func WriteJSON(w http.ResponseWriter, data any) {
	// First set the headers
	w.Header().Set("Content-Type", "application/json")

	// Write the data
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
