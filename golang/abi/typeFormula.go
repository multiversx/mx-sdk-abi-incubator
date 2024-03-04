package abi

import (
	"fmt"
	"strings"
)

type TypeFormula struct {
	Name           string
	TypeParameters []*TypeFormula
}

func NewTypeFormula(name string, typeParameters []*TypeFormula) *TypeFormula {
	return &TypeFormula{
		Name:           name,
		TypeParameters: typeParameters,
	}
}

func (t *TypeFormula) String() string {
	if len(t.TypeParameters) > 0 {
		var typeParameters []string
		for _, typeParameter := range t.TypeParameters {
			typeParameters = append(typeParameters, typeParameter.String())
		}
		return fmt.Sprintf("%s<%s>", t.Name, strings.Join(typeParameters, ", "))
	} else {
		return t.Name
	}
}
