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

package validation_test

import (
	"testing"

	"github.com/jonathanlloyd/rif/internal/app/fileversions"
	"github.com/jonathanlloyd/rif/internal/app/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidFileValidates(t *testing.T) {
	version := 0
	url := "http://example.com"
	method := "GET"

	yamlFileStruct := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        &url,
		Method:     &method,
	}

	errs := validation.ValidateRifYamlFile(yamlFileStruct)

	assert.Equal(t, 0, len(errs))
}

func TestRifFileVersionShouldBeRequired(t *testing.T) {
	url := "http://example.com"
	method := "GET"

	yamlFileStruct := fileversions.RifYamlFileV0{
		URL:    &url,
		Method: &method,
	}

	errs := validation.ValidateRifYamlFile(yamlFileStruct)

	assert.Equal(t, 1, len(errs))
	assert.Contains(t, "rif_version is required", errs[0].Error())
}
