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

package fileversions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jonathanlloyd/rif/internal/app/variables"
)

var httpMethods = map[string]bool{
	"CONNECT": true,
	"DELETE":  true,
	"GET":     true,
	"HEAD":    true,
	"OPTIONS": true,
	"PATCH":   true,
	"POST":    true,
	"PUT":     true,
	"TRACE":   true,
}

// RifFileV0 is the canonical in-memory representation of a parsed rif file
type RifFileV0 struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    *string
}

// RifYamlFileV0 is the struct used to unmarshal V0 of the RIF file format
type RifYamlFileV0 struct {
	RifVersion *int                         `yaml:"rif_version"`
	URL        *string                      `yaml:"url"`
	Method     *string                      `yaml:"method"`
	Headers    map[string]string            `yaml:"headers"`
	Body       string                       `yaml:"body"`
	Variables  map[string]RifYamlVariableV0 `yaml:"variables"`
}

// Canonicalize returns the parsed yaml in the canonical internal
// format, decoupling internal logic from the specifics of the yaml schema
// itself
func (y RifYamlFileV0) Canonicalize() (
	rFile RifFileV0,
	varDefs map[string]variables.VarDef,
	validationErrors []error,
) {
	canonicalURL, errs := canonicalizeURL(y.URL)
	if errs != nil {
		validationErrors = append(validationErrors, errs...)
	}

	canonicalMethod, errs := canonicalizeMethod(y.Method)
	if errs != nil {
		validationErrors = append(validationErrors, errs...)
	}

	rFile = RifFileV0{
		URL:     canonicalURL,
		Method:  canonicalMethod,
		Headers: y.Headers,
		Body:    &y.Body,
	}

	varDefs = map[string]variables.VarDef{}
	for varName, rawVarDef := range y.Variables {
		canonicalVarDef, err := rawVarDef.canonicalize()
		if err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("Variable \"%s\" is invalid: %s", varName, err.Error()),
			)
			continue
		}
		varDefs[varName] = canonicalVarDef
	}

	return rFile, varDefs, validationErrors
}

func canonicalizeURL(yamlURL *string) (
	canonicalURL string,
	validationErrors []error,
) {
	if yamlURL == nil {
		return "", []error{errors.New("Field \"URL\" is required")}
	}
	return *yamlURL, nil
}

func canonicalizeMethod(yamlMethod *string) (
	canonicalMethod string,
	validationErrors []error,
) {
	validationErrors = []error{}

	var methodMissing bool
	if yamlMethod == nil {
		validationErrors = append(
			validationErrors,
			errors.New("Field \"method\" is required"),
		)
		methodMissing = true
	} else {
		canonicalMethod = *yamlMethod
	}

	canonicalMethod = strings.ToUpper(canonicalMethod)

	isValidMethod := httpMethods[canonicalMethod]
	if !methodMissing && !isValidMethod {
		validationErrors = append(
			validationErrors,
			fmt.Errorf("Method \"%s\" is invalid", canonicalMethod),
		)
	}

	return canonicalMethod, validationErrors
}

// RifYamlVariableV0 is a struct used to unmarshal the variable schema portion
// of V0 of the RIF file format
type RifYamlVariableV0 struct {
	Type    string      `yaml:"type"`
	Default interface{} `yaml:"default"`
}

func (v RifYamlVariableV0) canonicalize() (variables.VarDef, error) {
	varType, err := v.canonicalizeType()
	if err != nil {
		return variables.VarDef{}, err
	}
	varDefault := v.canonicalizeDefault()

	return variables.VarDef{
		Type:    varType,
		Default: varDefault,
	}, nil
}

func (v RifYamlVariableV0) canonicalizeType() (variables.VarType, error) {
	var varType variables.VarType
	switch v.Type {
	case "boolean":
		varType = variables.Boolean
	case "number":
		varType = variables.Number
	case "string":
		varType = variables.String
	default:
		return varType, fmt.Errorf(
			"variable has invalid type \"%s\" "+
				"(valid types are boolean, number and string)",
			v.Type,
		)
	}

	return varType, nil
}

func (v RifYamlVariableV0) canonicalizeDefault() interface{} {
	var d interface{}

	switch value := v.Default.(type) {
	case int:
		d = int64(value)
	case int32:
		d = int64(value)
	case float32:
		d = float64(value)
	default:
		d = value
	}

	return d
}
