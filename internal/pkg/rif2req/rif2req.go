// RIF - (HTTP) Requests In Files

// Copyright (C) 2017 - Jonathan Lloyd (jonathan@thisisjonathan.com)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rif2req

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/jonathanlloyd/rif/internal/app/fileversions"
)

// Rif2Req takes a parsed .rif file and returns an equivalent stdlib http
// request struct
func Rif2Req(rFile fileversions.RifFileV0, rifVersion string) (*http.Request, error) {
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
