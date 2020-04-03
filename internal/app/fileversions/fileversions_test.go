package fileversions_test

import (
	"errors"
	"testing"

	"github.com/jonathanlloyd/rif/internal/app/fileversions"
	"github.com/jonathanlloyd/rif/internal/app/variables"
	"github.com/stretchr/testify/assert"
)

func TestValidFileValidates(t *testing.T) {
	version := 0
	URL := "http://example.com"
	method := "GET"

	validYamlFile := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        &URL,
		Method:     &method,
		Headers:    map[string]string{},
		Body:       "",
		Variables: map[string]fileversions.RifYamlVariableV0{
			"foo": {
				Type:    "string",
				Default: "bar",
			},
		},
	}

	rFile, varDefs, errs := validYamlFile.Canonicalize()

	emptyBody := ""
	expectedRFile := fileversions.RifFileV0{
		URL:     "http://example.com",
		Method:  "GET",
		Headers: map[string]string{},
		Body:    &emptyBody,
	}

	expectedVarDefs := map[string]variables.VarDef{
		"foo": {
			Type:    variables.String,
			Default: "bar",
		},
	}

	assert.Equal(t, expectedRFile, rFile)
	assert.Equal(t, expectedVarDefs, varDefs)
	assert.Empty(t, errs)
}

func TestLowercaseMethodsWork(t *testing.T) {
	version := 0
	URL := "http://example.com"
	method := "get"

	validYamlFile := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        &URL,
		Method:     &method,
		Headers:    map[string]string{},
		Body:       "",
		Variables: map[string]fileversions.RifYamlVariableV0{
			"foo": {
				Type:    "string",
				Default: "bar",
			},
		},
	}

	rFile, _, errs := validYamlFile.Canonicalize()

	emptyBody := ""
	expectedRFile := fileversions.RifFileV0{
		URL:     "http://example.com",
		Method:  "GET",
		Headers: map[string]string{},
		Body:    &emptyBody,
	}

	assert.Equal(t, expectedRFile, rFile)
	assert.Empty(t, errs)
}

func TestMissingURL(t *testing.T) {
	version := 0
	method := "GET"

	validYamlFile := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        nil,
		Method:     &method,
		Headers:    map[string]string{},
		Body:       "",
		Variables: map[string]fileversions.RifYamlVariableV0{
			"foo": {
				Type:    "string",
				Default: "bar",
			},
		},
	}

	_, _, errs := validYamlFile.Canonicalize()

	expectedErrs := []error{
		errors.New("Field \"url\" is required"),
	}

	assert.NotNil(t, errs)
	assert.ElementsMatch(t, expectedErrs, errs)
}

func TestMissingMethod(t *testing.T) {
	version := 0
	URL := "http://example.com"

	validYamlFile := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        &URL,
		Method:     nil,
		Headers:    map[string]string{},
		Body:       "",
		Variables: map[string]fileversions.RifYamlVariableV0{
			"foo": {
				Type:    "string",
				Default: "bar",
			},
		},
	}

	_, _, errs := validYamlFile.Canonicalize()

	expectedErrs := []error{
		errors.New("Field \"method\" is required"),
	}

	assert.NotNil(t, errs)
	assert.ElementsMatch(t, expectedErrs, errs)
}

func TestInvalidMethod(t *testing.T) {
	version := 0
	URL := "http://example.com"
	method := "FOOBAR"

	validYamlFile := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        &URL,
		Method:     &method,
		Headers:    map[string]string{},
		Body:       "",
		Variables: map[string]fileversions.RifYamlVariableV0{
			"foo": {
				Type:    "string",
				Default: "bar",
			},
		},
	}

	_, _, errs := validYamlFile.Canonicalize()

	expectedErrs := []error{
		errors.New("Method \"FOOBAR\" is invalid"),
	}

	assert.NotNil(t, errs)
	assert.ElementsMatch(t, expectedErrs, errs)
}

func TestInvalidVarDefType(t *testing.T) {
	version := 0
	URL := "http://example.com"
	method := "POST"

	validYamlFile := fileversions.RifYamlFileV0{
		RifVersion: &version,
		URL:        &URL,
		Method:     &method,
		Headers:    map[string]string{},
		Body:       "",
		Variables: map[string]fileversions.RifYamlVariableV0{
			"foo": {
				Type: "notavalidtype",
			},
		},
	}

	_, _, errs := validYamlFile.Canonicalize()

	expectedErrs := []error{
		errors.New(
			"Variable \"foo\" is invalid: " +
				"variable has invalid type \"notavalidtype\" " +
				"(valid types are boolean, number and string)",
		),
	}

	assert.NotNil(t, errs)
	assert.ElementsMatch(t, expectedErrs, errs)
}
