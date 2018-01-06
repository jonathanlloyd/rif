package rif2req_test

import (
	"github.com/turingincomplete/rif/internal/pkg/rif2req"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		req, err := rif2req.Rif2Req(rFile)

		assert.Nil(t, err)
		assert.Equal(t, req.Method, method)
	}
}

// Rif2Req should reject invalid HTTP methods
func TestInvalidMethods(t *testing.T) {
	rFile := rif2req.RifFileV0{
		URL:    "http://example.com/test",
		Method: "BAD_METHOD",
	}

	_, err := rif2req.Rif2Req(rFile)
	assert.NotNil(t, err)
}
