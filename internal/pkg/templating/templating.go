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

package templating

import "fmt"

// ApplyTemplate is used to substitute variable values into a given template
// string. If the template string is malformed or a variable is missing from
// the variable name -> value map, then an error is returned.
// Template Language:
//   The template language provided by this module looks like this:
//     "This is a variable: $(VAR_NAME)"
//   When { "VAR_NAME": "value" } is substituted into this template
//   ApplyTemplate will return:
//     "This is a variable: value"
func ApplyTemplate(templateStr string, vars map[string]string) (string, error) {
	l := lexer{input: templateStr}

	state := charState(&l)
	for state != nil {
		state = state(&l)
	}
	if l.errMsg != "" {
		return "", fmt.Errorf(l.errMsg)
	}

	var output string

	for i, tok := range l.tokens {
		isVar := i%2 != 0
		if isVar {
			value, ok := vars[tok]
			if !ok {
				return "", fmt.Errorf(
					"Template variable %s not found in data provided",
					tok,
				)
			}
			output += value
		} else {
			output += tok
		}
	}

	return output, nil
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input  string
	start  int
	pos    int
	tokens []string
	errMsg string
}

func charState(l *lexer) stateFn {
	for ; ; l.pos++ {
		if l.pos > len(l.input)-1 {
			l.tokens = append(l.tokens, l.input[l.start:l.pos])
			return nil
		}

		char := l.input[l.pos]
		var prevChar byte
		if l.pos != 0 {
			prevChar = l.input[l.pos-1]
		}

		if prevChar == '$' && char == '(' {
			l.tokens = append(l.tokens, l.input[l.start:l.pos-1])
			l.start, l.pos = l.pos+1, l.pos+1
			return varState
		}
	}
}

func varState(l *lexer) stateFn {
	for ; ; l.pos++ {
		if l.pos > len(l.input)-1 {
			l.errMsg = fmt.Sprintf("Template string terminates in a variable")
			return nil
		}

		char := l.input[l.pos]
		var prevChar byte
		if l.pos != 0 {
			prevChar = l.input[l.pos-1]
		}

		switch char {
		case '(':
			if prevChar == '$' {
				l.errMsg = fmt.Sprintf("Illegal nested variable at character %d", l.pos)
				return nil
			}
		case ')':
			l.tokens = append(l.tokens, l.input[l.start:l.pos])
			l.start, l.pos = l.pos+1, l.pos+1
			return charState
		default:
			continue
		}
	}
}
