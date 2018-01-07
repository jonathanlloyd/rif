package rif2req_test

import (
	"github.com/turingincomplete/rif/internal/pkg/rif2req"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Rif2Req should correctly set req method
func TestMethodSet(t *testing.T) {
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "GET",
	}
	req, _ := rif2req.Rif2Req(rFile, "1.0.0")

	assert.Equal(t, rFile.Method, req.Method)
}

// Rif2Req should correctly set req URL
func TestURLSet(t *testing.T) {
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "GET",
	}
	req, _ := rif2req.Rif2Req(rFile, "1.0.0")

	assert.Equal(t, "http", req.URL.Scheme)
	assert.Equal(t, "example.com", req.URL.Host)
	assert.Equal(t, "/test", req.URL.Path)
}

// Rif2Req should correctly set user agent
func TestUserAgentSet(t *testing.T) {
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "GET",
	}
	req, _ := rif2req.Rif2Req(rFile, "1.0.0")

	reqUserAgent := strings.Join(req.Header["User-Agent"], "")
	assert.Equal(t, "RIF/1.0.0", reqUserAgent)
}

// Rif2Req should correctly set additional req headers
func TestHeadersSet(t *testing.T) {
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "GET",
		Headers: map[string]string{
			"x-fake-header": "value",
		},
	}
	req, _ := rif2req.Rif2Req(rFile, "1.0.0")

	header := strings.Join(req.Header["X-Fake-Header"], "")
	assert.Equal(t, "value", header)
}

// Rif2Req should correctly set req body
func TestBodySet(t *testing.T) {
	body := "test_body"
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "POST",
		Body:   &body,
	}
	req, _ := rif2req.Rif2Req(rFile, "1.0.0")

	reqBody, err := ioutil.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, "test_body", string(reqBody))
}

// Rif2Req should accept valid HTTP methods
func TestValidMethods(t *testing.T) {
	httpMethods := []string{
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

	for _, method := range httpMethods {
		rFile := rif2req.RifFileV0{
			URL:    "http://example.com/test",
			Method: method,
		}
		_, err := rif2req.Rif2Req(rFile, "1.0.0")

		assert.Nil(t, err)
	}
}

// Rif2Req should reject invalid HTTP methods
func TestInvalidMethods(t *testing.T) {
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "BAD_METHOD",
	}

	_, err := rif2req.Rif2Req(rFile, "1.0.0")
	assert.NotNil(t, err)
}
