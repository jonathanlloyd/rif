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

	"gopkg.in/yaml.v2"
)

const usage = `RIF - (HTTP) Requests In Files

Usage:
  rif <filename>
  rif -h | --help
  rif --version

Options:
  -h --help     Show this screen.
  --version     Display the current version and build number.`

var (
	version string
	buildNo string
)

type RifFile struct {
	RifVersion int    `yaml:"rif_version"`
	URL        string `yaml:"url"`
	Method     string `yaml:"method"`
}

func main() {
	versionString := fmt.Sprintf("Version: %s\nBuild: %s", version, buildNo)
	userAgent := fmt.Sprintf("RIF/%s", version)
	arguments, _ := docopt.Parse(usage, nil, true, versionString, false)
	// Grab the name of the .rif file
	filenameInt, ok := arguments["<filename>"]
	if !ok {
		fmt.Println("Please specify a filename")
	}
	filename := filenameInt.(string)

	// Parse it
	rawFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .rif file: %s\n", err.Error())
		os.Exit(1)
	}

	rifFile := RifFile{}
	err = yaml.Unmarshal(rawFile, &rifFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing .rif file: %s\n", err.Error())
		os.Exit(1)
	}

	// Make the request
	client := &http.Client{}

	req, err := http.NewRequest(rifFile.Method, rifFile.URL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %s\n", err.Error())
		os.Exit(1)
	}
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %s\n", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(body))
}
