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

package templating_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turingincomplete/rif/internal/pkg/templating"
)

func TestParse(t *testing.T) {
	renderFunc, err := templating.Parse("The $(THING)'s name was $(NAME)")
	assert.Nil(t, err)

	output, err := renderFunc(map[string]string{
		"THING": "man",
		"NAME":  "Bob",
	})
	assert.Nil(t, err)

	assert.Equal(t, "The man's name was Bob", output)
}

func TestStringWithNoVariablesUnaltered(t *testing.T) {
	renderFunc, err := templating.Parse("No variables here")
	assert.Nil(t, err)

	output, err := renderFunc(map[string]string{})
	assert.Nil(t, err)

	assert.Equal(t, "No variables here", output)
}

func TestVariablesShouldBeClosed(t *testing.T) {
	testCases := []string{
		"$(a  $(b)",
		"$(a) $(b ",
	}

	for _, testCase := range testCases {
		_, err := templating.Parse(testCase)
		assert.NotNil(t, err)
	}
}

func TestVariablesShouldNotBeNested(t *testing.T) {
	testCases := []string{
		"$($(a))",
		"$($(a)",
	}

	for _, testCase := range testCases {
		_, err := templating.Parse(testCase)
		assert.NotNil(t, err)
	}
}

func TestAllTemplateVariablesRequired(t *testing.T) {
	renderFunc, err := templating.Parse("The $(THING)'s name was $(NAME)")
	assert.Nil(t, err)

	_, err = renderFunc(map[string]string{
		"THING": "man",
	})
	assert.NotNil(t, err)
}
