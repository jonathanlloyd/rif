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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/docopt/docopt-go"

	"github.com/turingincomplete/rif/internal/app/variables"
	"github.com/turingincomplete/rif/internal/pkg/rif2req"
	"github.com/turingincomplete/rif/internal/pkg/templating"
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

Example:
  rif my-request.rif count=12 secret=password
`

const majorVersion = 0

var (
	version string
	buildNo string
)

type rifYamlFile struct {
	RifVersion int               `yaml:"rif_version"`
	URL        string            `yaml:"url"`
	Method     string            `yaml:"method"`
	Headers    map[string]string `yaml:"headers"`
	Body       string            `yaml:"body"`
	Variables  map[string]struct {
		Type    string      `yaml:"type"`
		Default interface{} `yaml:"default"`
	} `yaml:"variables"`
}

// nolint: gocyclo
func main() {
	versionString := fmt.Sprintf("Version: %s\nBuild: %s", version, buildNo)
	arguments, _ := docopt.Parse(usage, nil, true, versionString, false)

	// Parse file
	filename := arguments["<filename>"].(string)

	rawFile, err := ioutil.ReadFile(filename)
	if err != nil {
		errorAndExit("Error reading .rif file", err)
	}

	rFile := rifYamlFile{}
	err = yaml.Unmarshal(rawFile, &rFile)
	if err != nil {
		errorAndExit("Error parsing .rif file", err)
	}

	if rFile.RifVersion > majorVersion {
		errorAndExit("Error parsing .rif file", fmt.Errorf(
			"rif file version greater than maxium supported version - %d",
			majorVersion))
	}

	// Work out variable values
	rawVars := arguments["<variable>"].([]string)
	inputVars, err := parseInputVars(rawVars)
	if err != nil {
		errorAndExit("Error parsing variables", err)
	}

	varDefinitions, err := preprocessVariableDefinitions(rFile)
	if err != nil {
		errorAndExit("Invalid variable definition", err)
	}

	varMap, err := variables.MakeMap(varDefinitions, inputVars)
	if err != nil {
		errorAndExit("Encountered an error calculating template variables", err)
	}

	// Apply template substitutions
	urlTemplate, err := templating.Parse(rFile.URL)
	if err != nil {
		errorAndExit("Error reading URL template", err)
	}
	rFile.URL, err = urlTemplate(varMap)
	if err != nil {
		errorAndExit("Error rendering URL template", err)
	}

	newHeaders := map[string]string{}
	for headerName, headerValue := range rFile.Headers {
		headerTemplate, err := templating.Parse(headerValue)
		if err != nil {
			errorAndExit("Error reading header template", err)
		}
		renderedHeader, err := headerTemplate(varMap)
		if err != nil {
			errorAndExit("Error rendering header template", err)
		}
		newHeaders[headerName] = renderedHeader
	}
	rFile.Headers = newHeaders

	bodyTemplate, err := templating.Parse(rFile.Body)
	if err != nil {
		errorAndExit("Error reading body template", err)
	}
	rFile.Body, err = bodyTemplate(varMap)
	if err != nil {
		errorAndExit("Error rendering body template", err)
	}

	// Make the request
	req, err := rif2req.Rif2Req(
		rif2req.RifFileV0{
			URL:     rFile.URL,
			Method:  rFile.Method,
			Headers: rFile.Headers,
			Body:    &rFile.Body,
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
	defaultFormat := !ok
	httpFormat := outputFormat == "http" || outputFormat == "HTTP"

	if httpFormat {
		newReq, err := rif2req.Rif2Req(
			rif2req.RifFileV0{
				URL:     rFile.URL,
				Method:  rFile.Method,
				Headers: rFile.Headers,
				Body:    &rFile.Body,
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

func errorAndExit(errPrefix string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", errPrefix, err.Error())
	os.Exit(1)
}

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

func preprocessVariableDefinitions(
	rFile rifYamlFile,
) (map[string]variables.VarDef, error) {
	varDefinitions := map[string]variables.VarDef{}
	for varName, rawVarDef := range rFile.Variables {
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
				"variable definition \"%s\" has invalid type: \"%s\"",
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
