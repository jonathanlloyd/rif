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
//  - Parse the given RIF file to produce the request and variable definitions
//  - Parse the arguments given on the command-line and combine them with the
//    given variable definitions to produce a map from variable name -> value
//  - Interpolate any variables passed in into any templated fields in the
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
	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], versionString)

	/**
	 * Step 1 - Parse the given RIF file to produce a request definition
	 * and a map defining the schema for the variables declared in the file
	**/
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

	reqDef, varDefs, canonicalizationErrors := yamlStruct.Canonicalize()
	validationErrs = append(validationErrs, canonicalizationErrors...)

	if len(validationErrs) > 0 {
		errString := ""
		for _, err := range validationErrs {
			errString += "\n - " + err.Error()
		}
		errorAndExit("Invalid .rif file", errors.New(errString))
	}

	/**
	 * Step 2 - Parse the variables given over the command line and combine
	 * them with the default values in the RIF file's variable definitions to
	 * produce a mapping from variable name to value.
	**/
	rawVars := arguments["<variable>"].([]string)

	inputVarMap, err := parseInputVars(rawVars)
	if err != nil {
		errorAndExit("Invalid parameters", err)
	}

	varMap, err := calculateVariableValues(inputVarMap, varDefs)
	if err != nil {
		errorAndExit("Invalid parameters", err)
	}

	/**
	 * Step 3 - Substitute the calculated variable values into any template
	 * strings in the request definition
	**/
	reqDef, err = substituteVariableValues(varMap, reqDef)
	if err != nil {
		errorAndExit("Error interpolating variables", err)
	}

	/**
	 * Step 4 - Execute the http request
	**/
	req, err := rif2req.Rif2Req(reqDef, version)
	if err != nil {
		errorAndExit("Error making request", err)
	}
	// Need a second copy for printing in different output modes
	// Performing the request consumes the body reader
	printReq, _ := rif2req.Rif2Req(reqDef, version)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorAndExit("Error making request", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	/**
	 * Step 5 - Print the request/response in the selected format
	**/
	outputFormat, ok := arguments["--output"].(string)
	defaultFormat := !ok || outputFormat == "default"
	httpFormat := outputFormat == "http" || outputFormat == "HTTP"
	curlFormat := outputFormat == "curl" || outputFormat == "cURL"

	if httpFormat {
		httpReq, err := httputil.DumpRequestOut(printReq, true)
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
		curlOutput, err := http2curl.GetCurlCommand(printReq)
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
	inputVarMap map[string]string,
	varDefs map[string]variables.VarDef,
) (map[string]string, error) {
	emptyMap := map[string]string{}

	err := variables.ValidateInputVars(varDefs, inputVarMap)
	if err != nil {
		return emptyMap, err
	}

	varMap, err := variables.MakeMap(varDefs, inputVarMap)
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

// substituteVariableValues takes a map from variable name to value and a
// parsed RIF file and interpolates the given variables into any template
// strings that exist in the file.
func substituteVariableValues(
	varMap map[string]string,
	reqDef fileversions.RifFileV0,
) (fileversions.RifFileV0, error) {
	emptyReqDef := fileversions.RifFileV0{}

	renderedURL, urlErr := templating.ApplyTemplate(reqDef.URL, varMap)
	if urlErr != nil {
		return emptyReqDef, urlErr
	}
	reqDef.URL = renderedURL

	newHeaders := map[string]string{}
	for headerName, headerValue := range reqDef.Headers {
		renderedHeader, headerErr := templating.ApplyTemplate(headerValue, varMap)
		if headerErr != nil {
			return emptyReqDef, headerErr
		}
		newHeaders[headerName] = renderedHeader
	}
	reqDef.Headers = newHeaders

	renderedBody, bodyErr := templating.ApplyTemplate(*reqDef.Body, varMap)
	if bodyErr != nil {
		return emptyReqDef, bodyErr
	}
	reqDef.Body = &renderedBody

	return reqDef, nil
}
