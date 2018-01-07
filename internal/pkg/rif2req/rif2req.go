package rif2req

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// RifFileV0 is an in-memory representation of the unversioned beta .rif file
// format
type RifFileV0 struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    *string
}

var httpMethods = []string{
	"CONNECT",
	"DELETE",
	"GET",
	"HEAD",
	"OPTIONS",
	"PATCH",
	"POST",
	"PUT",
	"TRACE",
}

// Rif2Req takes a parsed .rif file and returns an equivalent stdlib http
// request struct
func Rif2Req(rFile RifFileV0, rifVersion string) (*http.Request, error) {
	// Validate rFile
	isValidMethod := false
	for _, method := range httpMethods {
		if method == rFile.Method {
			isValidMethod = true
			break
		}
	}
	if !isValidMethod {
		return nil, fmt.Errorf("Method %s is invalid", rFile.Method)
	}

	// Create request
	var body io.Reader
	if rFile.Body != nil {
		body = bytes.NewBuffer([]byte(*rFile.Body))
	}
	req, err := http.NewRequest(rFile.Method, rFile.URL, body)
	if err != nil {
		return nil, err
	}

	// Add additional headers
	for headerName, headerValue := range rFile.Headers {
		req.Header.Add(headerName, headerValue)
	}

	// Add user agent
	userAgent := fmt.Sprintf("RIF/%s", rifVersion)
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}
