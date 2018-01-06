package rif2req

import (
	"fmt"
	"net/http"
)

type RifFileV0 struct {
	URL    string
	Method string
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

func Rif2Req(rFile RifFileV0) (*http.Request, error) {
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

	req, err := http.NewRequest(rFile.Method, rFile.URL, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
