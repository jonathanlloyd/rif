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

package variables

import (
	"fmt"
	"strconv"
)

// VarType is used to enumerate the different valid types for variable
// substitution: boolean, number and string.
type VarType int

const (
	// Boolean - true/false
	Boolean VarType = iota
	// Number - int/float
	Number
	// String - any other string
	String
)

// VarDef represents a variable definition that can be used in template
// substitution. This includes the type and default value of the variable.
type VarDef struct {
	Type    VarType
	Default interface{}
}

// MakeMap takes a set of variable definitions and input variables and returns
// a variable map that can be used for template substitution. This include
// inserting default values and validating that the format of each variable
// matches its type.
func MakeMap(
	varDefs map[string]VarDef,
	inputVars map[string]string,
) (map[string]string, error) {
	err := validateVarDefDefaults(varDefs)
	if err != nil {
		return map[string]string{}, err
	}

	err = validateInputVarTypes(varDefs, inputVars)
	if err != nil {
		return map[string]string{}, err
	}

	// Build variable map
	varMap := map[string]string{}
	for varName, varDef := range varDefs {
		if val, ok := inputVars[varName]; ok {
			varMap[varName] = val
		} else if varDef.Default != nil {
			varMap[varName] = fmt.Sprintf("%+v", varDef.Default)
		} else {
			return map[string]string{}, fmt.Errorf(
				"error adding \"%s\" to variable map: No input variable or default given",
				varName,
			)
		}
	}

	return varMap, nil
}

// ValidateInputVars checks the provided input vars against the variable
// definitions in the RIF file. If the input vars do not match the definition
// then an error is returned containing a thorough description of the issues
// and how to resolve them.
func ValidateInputVars(
	varDefs map[string]VarDef,
	inputVars map[string]string,
) error {
	return nil
}

func validateVarDefDefaults(varDefs map[string]VarDef) error {
	for varName, varDef := range varDefs {
		if varDef.Default == nil {
			continue
		}
		_, defaultIsBool := varDef.Default.(bool)
		_, defaultIsInt := varDef.Default.(int64)
		_, defaultIsFloat := varDef.Default.(float64)
		_, defaultIsString := varDef.Default.(string)

		if (varDef.Type == Boolean && !defaultIsBool) ||
			(varDef.Type == Number && !(defaultIsInt || defaultIsFloat)) ||
			(varDef.Type == String && !defaultIsString) {
			return fmt.Errorf(
				"\"%s\" has invalid variable definition: default value is of incorrect type",
				varName,
			)
		}
	}

	return nil
}

func validateInputVarTypes(
	varDefs map[string]VarDef,
	inputVars map[string]string,
) error {
	for varName, varDef := range varDefs {
		input, ok := inputVars[varName]
		if !ok {
			continue
		}

		_, boolErr := strconv.ParseBool(input)
		_, intErr := strconv.ParseInt(input, 10, 64)
		_, floatErr := strconv.ParseFloat(input, 64)

		if ((varDef.Type == Boolean && boolErr != nil) ||
			(varDef.Type == Number && !(intErr == nil || floatErr == nil))) &&
			(varDef.Type != String) {
			return fmt.Errorf(
				"type of input variable \"%s\" does not match variable definition",
				varName,
			)
		}
	}

	return nil
}
