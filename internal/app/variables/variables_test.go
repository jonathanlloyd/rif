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
	"testing"

	"github.com/jonathanlloyd/rif/internal/app/variables"
	"github.com/stretchr/testify/assert"
)

func TestMakeMapShouldWorkForAllTypes(t *testing.T) {
	variableDefinitions := []map[string]variables.VarDef{
		{
			"VAR": variables.VarDef{
				Type: variables.Boolean,
			},
		},
		{
			"VAR": variables.VarDef{
				Type: variables.Number,
			},
		},
		{
			"VAR": variables.VarDef{
				Type: variables.String,
			},
		},
	}

	inputVariables := []map[string]string{
		{
			"VAR": "true",
		},
		{
			"VAR": "12.3",
		},
		{
			"VAR": "some_string",
		},
	}

	for i := range variableDefinitions {
		varDef := variableDefinitions[i]
		inputVar := inputVariables[i]

		variableMap, err := variables.MakeMap(varDef, inputVar)
		assert.Nil(t, err)

		finalVar := variableMap["VAR"]
		expectedValue := inputVar["VAR"]

		assert.Equal(t, expectedValue, finalVar)
	}
}

func TestMakeMapShouldOverrideDefaults(t *testing.T) {
	variableDefinition := map[string]variables.VarDef{
		"DONT_OVERRIDE_ME": {
			Type:    variables.Boolean,
			Default: true,
		},
		"OVERRIDE_ME": {
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
		"NO_VALUE": {
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

	for i := range variableDefinitions {
		varDefs := variableDefinitions[i]
		inputVars := inputVariables[i]
		_, err := variables.MakeMap(varDefs, inputVars)
		assert.NotNil(t, err)
	}
}

func TestValidateInputVarsShouldAcceptValidInput(t *testing.T) {
	varDefs := map[string]variables.VarDef{
		"REQUIRED_1": variables.VarDef{
			Type: variables.String,
		},
		"REQUIRED_2": variables.VarDef{
			Type: variables.String,
		},
		"OPTIONAL_1": variables.VarDef{
			Type:    variables.String,
			Default: "VALUE",
		},
		"OPTIONAL_2": variables.VarDef{
			Type:    variables.String,
			Default: "VALUE",
		},
	}

	inputVars := map[string]string{
		"REQUIRED_1": "VALUE",
		"REQUIRED_2": "VALUE",
	}

	err := variables.ValidateInputVars(varDefs, inputVars)

	assert.Nil(t, err)
}

func TestValidateInputVarsShouldReturnNiceError(t *testing.T) {
	varDefs := map[string]variables.VarDef{
		"REQUIRED_1": variables.VarDef{
			Type: variables.String,
		},
		"REQUIRED_2": variables.VarDef{
			Type: variables.String,
		},
		"OPTIONAL_1": variables.VarDef{
			Type:    variables.String,
			Default: "VALUE",
		},
		"OPTIONAL_2": variables.VarDef{
			Type:    variables.String,
			Default: "VALUE",
		},
	}

	inputVars := map[string]string{
		"REQUIRED_1": "VALUE",
	}

	expectedErrorMessage := `
Missing required variable(s): REQUIRED_2

The variables for this RIF file are as follows:
Required:
 - REQUIRED_1 ( string )
 - REQUIRED_2 ( string )
Optional:
 - OPTIONAL_1 ( string, default=VALUE )
 - OPTIONAL_2 ( string, default=VALUE )
`

	err := variables.ValidateInputVars(varDefs, inputVars)
	assert.Equal(t, expectedErrorMessage, err.Error())
}
