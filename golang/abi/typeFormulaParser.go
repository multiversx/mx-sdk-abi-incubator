package abi

import (
	"errors"
	"fmt"
	"strings"
)

type TypeFormulaParser struct {
	BEGIN_TYPE_PARAMETERS string
	END_TYPE_PARAMETERS   string
	COMMA                 string
	PUNCTUATION           []string
}

func NewTypeFormulaParser() *TypeFormulaParser {
	return &TypeFormulaParser{
		BEGIN_TYPE_PARAMETERS: "<",
		END_TYPE_PARAMETERS:   ">",
		COMMA:                 ",",
		PUNCTUATION:           []string{",", "<", ">"},
	}
}

func (p *TypeFormulaParser) ParseExpression(expression string) (*TypeFormula, error) {
	expression = strings.TrimSpace(expression)

	tokens := p.tokenizeExpression(expression)
	tokens = filter(tokens, func(token string) bool { return token != p.COMMA })

	stack := make([]any, 0)

	for _, token := range tokens {
		if contains(p.PUNCTUATION, token) {
			if token == p.END_TYPE_PARAMETERS {
				typeParameters := make([]*TypeFormula, 0)

				for {
					if len(stack) == 0 {
						return nil, errors.New("Badly specified type parameters.")
					}

					if stack[len(stack)-1] == p.BEGIN_TYPE_PARAMETERS {
						break
					}

					item := stack[len(stack)-1]
					if typeFormula, ok := item.(*TypeFormula); ok {
						typeParameters = append(typeParameters, typeFormula)
					} else {
						typeFormula = NewTypeFormula(item.(string), []*TypeFormula{})
						typeParameters = append(typeParameters, typeFormula)
					}

					stack = stack[:len(stack)-1]
				}

				stack = stack[:len(stack)-1] // pop "<" symbol
				typeName := stack[len(stack)-1].(string)
				stack = stack[:len(stack)-1]
				typeFormula := NewTypeFormula(typeName, reverse(typeParameters))
				stack = append(stack, typeFormula)
			} else if token == p.BEGIN_TYPE_PARAMETERS {
				// The symbol is pushed as a simple string,
				// as it will never be interpreted, anyway.
				stack = append(stack, token)
			} else if token == p.COMMA {
				// We simply ignore commas
			} else {
				return nil, errors.New("Unexpected token (punctuation): " + token)
			}
		} else {
			// It's a type name. We push it as a simple string.
			stack = append(stack, token)
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("Unexpected stack length at end of parsing: %d", len(stack))
	}
	if _, ok := stack[0].(string); ok && contains(p.PUNCTUATION, stack[0].(string)) {
		return nil, errors.New("Unexpected root element.")
	}

	if typeName, ok := stack[0].(string); ok {
		// Expression contained a simple, non-generic type
		return NewTypeFormula(typeName, []*TypeFormula{}), nil
	} else if typeFormula, ok := stack[0].(*TypeFormula); ok {
		return typeFormula, nil
	} else {
		return nil, errors.New("Unexpected item on stack: " + stack[0].(string))
	}
}

func (p *TypeFormulaParser) tokenizeExpression(expression string) []string {
	tokens := make([]string, 0)
	currentToken := ""

	for i := 0; i < len(expression); i++ {
		character := string(expression[i])

		if !contains(p.PUNCTUATION, character) {
			// Non-punctuation character
			currentToken += character
		} else {
			if currentToken != "" {
				tokens = append(tokens, strings.TrimSpace(currentToken))
				currentToken = ""
			}
			// Punctuation character
			tokens = append(tokens, character)
		}
	}

	if currentToken != "" {
		tokens = append(tokens, strings.TrimSpace(currentToken))
	}

	return tokens
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func contains(vs []string, t string) bool {
	for _, v := range vs {
		if v == t {
			return true
		}
	}
	return false
}

// See: https://github.com/golang/go/wiki/SliceTricks#reversing
func reverse(vs []*TypeFormula) []*TypeFormula {
	for i := len(vs)/2 - 1; i >= 0; i-- {
		opp := len(vs) - 1 - i
		vs[i], vs[opp] = vs[opp], vs[i]
	}
	return vs
}
