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

package variables_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/turingincomplete/rif/internal/app/variables"
	"testing"
)

func TestMakeMapShouldOverrideDefaults(t *testing.T) {
	variableDefinition := map[string]variables.VarDef{
		"DONT_OVERRIDE_ME": variables.VarDef{
			Type:    variables.Boolean,
			Default: true,
		},
		"OVERRIDE_ME": variables.VarDef{
			Type:    variables.Boolean,
			Default: false,
		},
	}
	inputVariables := map[string]string{
		"OVERRIDE_ME": "true",
	}

	variableMap, err := variables.MakeMap(variableDefinition, inputVariables)
	assert.Nil(t, err)

	dontOverrideMe := variableMap["DONT_OVERRIDE_ME"]
	overrideMe := variableMap["OVERRIDE_ME"]

	assert.Equal(t, "true", dontOverrideMe)
	assert.Equal(t, "true", overrideMe)
}

func TestMakeMapShouldErrorIfNoValue(t *testing.T) {
	variableDefinition := map[string]variables.VarDef{
		"NO_VALUE": variables.VarDef{
			Type: variables.String,
		},
	}
	inputVariables := map[string]string{}

	_, err := variables.MakeMap(variableDefinition, inputVariables)
	assert.NotNil(t, err)
}

func TestMakeMapShouldErrorBadDefaultType(t *testing.T) {
	variableDefinitions := []map[string]variables.VarDef{
		{
			"BAD_DEFAULT": variables.VarDef{
				Type:    variables.Boolean,
				Default: "NOT_A_BOOL",
			},
		},
		{
			"BAD_DEFAULT": variables.VarDef{
				Type:    variables.Number,
				Default: "NOT_A_NUMBER",
			},
		},
	}
	inputVariables := map[string]string{}

	for _, varDef := range variableDefinitions {
		_, err := variables.MakeMap(varDef, inputVariables)
		assert.NotNil(t, err)
	}
}

func TestMakeMapShouldErrorBadInputType(t *testing.T) {
	variableDefinitions := []map[string]variables.VarDef{
		{
			"BAD_INPUT": variables.VarDef{
				Type: variables.Boolean,
			},
		},
		{
			"BAD_INPUT": variables.VarDef{
				Type: variables.Number,
			},
		},
	}
	inputVariables := []map[string]string{
		{
			"BAD_INPUT": "NOT_A_BOOL",
		},
		{
			"BAD_INPUT": "NOT_A_NUMBER",
		},
	}

	for i, _ := range variableDefinitions {
		varDefs := variableDefinitions[i]
		inputVars := inputVariables[i]
		_, err := variables.MakeMap(varDefs, inputVars)
		assert.NotNil(t, err)
	}
}
