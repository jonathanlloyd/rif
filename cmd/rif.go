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
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
  rif <filename> [<variable>]...
  rif -h | --help
  rif --version

Options:
  -h --help     Show this screen.
  --version     Display the current version and build number.

Example:
  rif my-request.rif count=12 secret=password
`

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
		fmt.Fprintf(os.Stderr, "Error reading .rif file: %s\n", err.Error())
		os.Exit(1)
	}

	rFile := rifYamlFile{}
	err = yaml.Unmarshal(rawFile, &rFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing .rif file: %s\n", err.Error())
		os.Exit(1)
	}

	// Work out variable values
	rawVars := arguments["<variable>"].([]string)
	inputVars, err := parseInputVars(rawVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	varDefinitions, err := preprocessVariableDefinitions(rFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	varMap, err := variables.MakeMap(varDefinitions, inputVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error calculating template variables: %s\n", err.Error())
		os.Exit(1)
	}

	// Apply template substitutions
	urlTemplate, err := templating.Parse(rFile.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading URL template: %s", err.Error())
		os.Exit(1)
	}
	rFile.URL, err = urlTemplate(varMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering URL template: %s", err.Error())
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error making request: %s\n", err.Error())
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(body))
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
