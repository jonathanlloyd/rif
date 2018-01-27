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

// RenderFunc is a function that represents a parsed template string.
// It takes a map which contains the variables to be substituted into
// the template. If a variable in the template string is not given
// an error will be returned.
type RenderFunc func(map[string]string) (string, error)

// Parse is a function that parses the given template string and returns
// a function that can be used to render that string with substitutions.
func Parse(templateStr string) (RenderFunc, error) {
	l := lexer{input: templateStr}

	state := charState(&l)
	for state != nil {
		state = state(&l)
	}
	if l.errMsg != "" {
		return nil, fmt.Errorf(l.errMsg)
	}

	render := func(data map[string]string) (string, error) {
		var output string

		for i, tok := range l.tokens {
			isVar := i%2 != 0
			if isVar {
				value, ok := data[tok]
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

	return render, nil
}
