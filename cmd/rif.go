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

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/docopt/docopt-go"

	"github.com/jonathanlloyd/rif/internal/app/fileversions"
	"github.com/jonathanlloyd/rif/internal/app/validation"
	"github.com/jonathanlloyd/rif/internal/app/variables"
	"github.com/jonathanlloyd/rif/internal/pkg/rif2req"
	"github.com/jonathanlloyd/rif/internal/pkg/templating"
	"github.com/moul/http2curl"
	"gopkg.in/yaml.v2"
)

const usage = `
                   
██████╗ ██╗███████╗
██╔══██╗██║██╔════╝
██████╔╝██║█████╗  
██╔══██╗██║██╔══╝  
██║  ██║██║██║     
╚═╝  ╚═╝╚═╝╚═╝     
                   
(HTTP) Requests In Files

Usage:
  rif <filename> [--output=<output-format>] [<variable>]...
  rif -h | --help
  rif --version

Options:
  -h --help     Show this screen.
  --version     Display the current version and build number.

Output Formats:
 - default: Returns the http response as given by the server
 - http: Returns a readable version of the raw HTTP request/response cycle
 - curl: Returns a cURL command equivalent to the given request

Example:
  rif my-request.rif count=12 secret=password --output=http
`

const majorVersion = 0

var (
	version string
	buildNo string
)

// RIF only really does one thing (make a HTTP request based on the file you
// give it), so this main function is really all there is to it.
// It should go something like this:
//  - Parse the arguments given on the command-line
//  - Parse the given RIF file
//  - Validate the given RIF file
//  - Iterpolate any variables passed in into any templated fields in the
//    RIF file. This can include any of the major parts of the request:
//      - URL
//      - Headers
//      - Body
//  - This provides us with everything we need to perform a HTTP request. What
//    we do next depends on the given output format:
//      - default: Do the request and print the response
//      - http: Do the request and print the full request/response cycle
//      - curl: Print a cURL command that is equivalent to the given request

// nolint: gocyclo
func main() {
	versionString := fmt.Sprintf("Version: %s\nBuild: %s", version, buildNo)
	arguments, _ := docopt.Parse(usage, nil, true, versionString, false)

	filename := arguments["<filename>"].(string)
	rawFile, err := ioutil.ReadFile(filename)
	if err != nil {
		errorAndExit("Error reading .rif file", err)
	}

	yamlStruct := fileversions.RifYamlFileV0{}
	err = yaml.Unmarshal(rawFile, &yamlStruct)
	if err != nil {
		errorAndExit("Error parsing .rif file", err)
	}

	validationErrs := validation.ValidateRifYamlFile(yamlStruct, majorVersion)
	if len(validationErrs) > 0 {
		errString := ""
		for _, err := range validationErrs {
			errString += "\n - " + err.Error()
		}
		errorAndExit("Invalid .rif file", errors.New(errString))
	}

	// Here we have an array of raw variables from the command line in the form:
	// "var_name=value"
	rawVars := arguments["<variable>"].([]string)
	// We need to parse the raw variables and combine them with the
	// schema defined in the yaml file to product a mapping from variable name to
	// value with all the defaulting and type checking sorted out.
	varMap, err := calculateVariableValues(rawVars, yamlStruct.Variables)
	if err != nil {
		errorAndExit("Invalid parameters", err)
	}

	// We can then iterpolate the variables into the different parts of the
	// request:
	yamlStruct, err = substituteVariableValues(varMap, yamlStruct)
	if err != nil {
		errorAndExit("Error interpolating variables", err)
	}

	// Make the request
	req, err := rif2req.Rif2Req(
		fileversions.RifFileV0{
			URL:     *yamlStruct.URL,
			Method:  *yamlStruct.Method,
			Headers: yamlStruct.Headers,
			Body:    &yamlStruct.Body,
		},
		version,
	)
	if err != nil {
		errorAndExit("Error making request", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorAndExit("Error making request", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Print the request/response in the appropriate format
	outputFormat, ok := arguments["--output"].(string)
	defaultFormat := !ok || outputFormat == "default"
	httpFormat := outputFormat == "http" || outputFormat == "HTTP"
	curlFormat := outputFormat == "curl" || outputFormat == "cURL"

	if httpFormat {
		newReq, err := rif2req.Rif2Req(
			fileversions.RifFileV0{
				URL:     *yamlStruct.URL,
				Method:  *yamlStruct.Method,
				Headers: yamlStruct.Headers,
				Body:    &yamlStruct.Body,
			},
			version,
		)
		if err != nil {
			// This shouldn't happen because we have already made a request
			// from this rif file, but I'm leaving this here for completeness.
			errorAndExit("Error printing HTTP request", err)
		}
		httpReq, err := httputil.DumpRequestOut(newReq, true)
		if err != nil {
			errorAndExit("Error printing HTTP request", err)
		}
		httpResp, err := httputil.DumpResponse(resp, true)
		if err != nil {
			errorAndExit("Error printing HTTP response", err)
		}
		fmt.Println("Request\n-------")
		fmt.Println(string(httpReq) + "\n")
		fmt.Println("Response\n--------")
		fmt.Println(string(httpResp))
	} else if curlFormat {
		newReq, err := rif2req.Rif2Req(
			fileversions.RifFileV0{
				URL:     *yamlStruct.URL,
				Method:  *yamlStruct.Method,
				Headers: yamlStruct.Headers,
				Body:    &yamlStruct.Body,
			},
			version,
		)
		if err != nil {
			// This shouldn't happen because we have already made a request
			// from this rif file, but I'm leaving this here for completeness.
			errorAndExit("Error printing cURL request", err)
		}

		curlOutput, err := http2curl.GetCurlCommand(newReq)
		if err != nil {
			errorAndExit("Error printing cURL request", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorAndExit("Error making request", err)
		}

		fmt.Println("cURL command")
		fmt.Println("------------")
		fmt.Println(curlOutput)
		fmt.Println("")
		fmt.Println("Response")
		fmt.Println("--------")
		fmt.Println(string(body))

	} else if defaultFormat {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorAndExit("Error making request", err)
		}
		fmt.Println(string(body))
	} else {
		errorAndExit("Unknown output format", fmt.Errorf(outputFormat))
	}
}

// errorAndExit takes a prefix message and an error, prints a pretty error
// message to the screen and exits the program returning an error code.
func errorAndExit(errPrefix string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", errPrefix, err.Error())
	os.Exit(1)
}

// calculateVariableValues takes the raw variable assignments from the
// command-line and the schema from the yaml file and returns a map from
// variable name to value with all the defaulting and type-checking sorted
// out.
func calculateVariableValues(
	rawVars []string,
	yamlVarDefs map[string]fileversions.RifYamlVariableV0,
) (map[string]string, error) {
	emptyMap := map[string]string{}

	inputVars, err := parseInputVars(rawVars)
	if err != nil {
		return emptyMap, err
	}

	variableSchema, err := makeVariableSchema(yamlVarDefs)
	if err != nil {
		return emptyMap, err
	}

	err = variables.ValidateInputVars(variableSchema, inputVars)
	if err != nil {
		return emptyMap, err
	}

	varMap, err := variables.MakeMap(variableSchema, inputVars)
	if err != nil {
		return emptyMap, err
	}

	return varMap, nil
}

// parseInputVars takes the raw variable assignments from the command-line
// and parses them into a map from variable name to value.
func parseInputVars(rawVars []string) (map[string]string, error) {
	vars := map[string]string{}

	for _, rawVar := range rawVars {
		parts := strings.Split(rawVar, "=")
		if len(parts) == 1 {
			return map[string]string{}, fmt.Errorf(
				"error parsing variable \"%s\": Variable names and values must be separated by \"=\"",
				rawVar,
			)
		}
		varName := parts[0]
		varValue := strings.Join(parts[1:], "=")

		vars[varName] = varValue
	}

	return vars, nil
}

// makeVariableSchema takes the variables section of the struct used to parse
// the yaml file and returns a decorated version where strings have been
// replaced with enums etc.
func makeVariableSchema(
	yamlVarDefs map[string]fileversions.RifYamlVariableV0,
) (map[string]variables.VarDef, error) {
	varDefinitions := map[string]variables.VarDef{}
	for varName, rawVarDef := range yamlVarDefs {
		var varType variables.VarType
		switch rawVarDef.Type {
		case "boolean":
			varType = variables.Boolean
		case "number":
			varType = variables.Number
		case "string":
			varType = variables.String
		default:
			return map[string]variables.VarDef{}, fmt.Errorf(
				"variable definition \"%s\" has invalid type \"%s\". "+
					"Valid types are boolean, number and string",
				varName,
				rawVarDef.Type,
			)
		}

		var varDefault interface{}
		switch value := rawVarDef.Default.(type) {
		case int:
			varDefault = int64(value)
		case int32:
			varDefault = int64(value)
		case float32:
			varDefault = float64(value)
		default:
			varDefault = value
		}
		varDefinitions[varName] = variables.VarDef{
			Type:    varType,
			Default: varDefault,
		}
	}

	return varDefinitions, nil
}

// substituteVariableValues takes a map from variable name to value and a
// parsed RIF file and interpolates the given variables into any template
// strings that exist in the file.
func substituteVariableValues(
	varMap map[string]string,
	yamlStruct fileversions.RifYamlFileV0,
) (fileversions.RifYamlFileV0, error) {
	applyTemplate := func(templateString string) (string, error) {
		template, err := templating.Parse(templateString)
		if err != nil {
			return "", err
		}
		renderedString, err := template(varMap)
		if err != nil {
			return "", err
		}
		return renderedString, nil
	}
	emptyYaml := fileversions.RifYamlFileV0{}

	renderedURL, urlErr := applyTemplate(*yamlStruct.URL)
	if urlErr != nil {
		return emptyYaml, urlErr
	}
	yamlStruct.URL = &renderedURL

	newHeaders := map[string]string{}
	for headerName, headerValue := range yamlStruct.Headers {
		renderedHeader, headerErr := applyTemplate(headerValue)
		if headerErr != nil {
			return emptyYaml, headerErr
		}
		newHeaders[headerName] = renderedHeader
	}
	yamlStruct.Headers = newHeaders

	renderedBody, bodyErr := applyTemplate(yamlStruct.Body)
	if bodyErr != nil {
		return emptyYaml, bodyErr
	}
	yamlStruct.Body = renderedBody

	return yamlStruct, nil
}
