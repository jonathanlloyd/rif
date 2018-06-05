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

// RifYamlFileV0 is the struct used to unmarshal V0 of the RIF file format
type RifYamlFileV0 struct {
	RifVersion *int                         `yaml:"rif_version"`
	URL        *string                      `yaml:"url"`
	Method     *string                      `yaml:"method"`
	Headers    map[string]string            `yaml:"headers"`
	Body       string                       `yaml:"body"`
	Variables  map[string]RifYamlVariableV0 `yaml:"variables"`
}

// RifYamlVariableV0 is a struct used to unmarshal the variable schema portion
// of V0 of the RIF file format
type RifYamlVariableV0 struct {
	Type    string      `yaml:"type"`
	Default interface{} `yaml:"default"`
}

// RifFileV0 is an in-memory representation of the unversioned beta .rif file
// format
type RifFileV0 struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    *string
}
