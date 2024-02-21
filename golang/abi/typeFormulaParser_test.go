package abi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeFormulaParser_ParseExpression(t *testing.T) {
	parser := NewTypeFormulaParser()

	formula, err := parser.ParseExpression("i64")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "i64", TypeParameters: []*TypeFormula{}}, *formula)
	require.Equal(t, "i64", formula.String())

	formula, err = parser.ParseExpression("  i64  ")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "i64", TypeParameters: []*TypeFormula{}}, *formula)
	require.Equal(t, "i64", formula.String())

	formula, err = parser.ParseExpression("utf-8 string")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "utf-8 string", TypeParameters: []*TypeFormula{}}, *formula)
	require.Equal(t, "utf-8 string", formula.String())

	formula, err = parser.ParseExpression("MultiResultVec<MultiResult2<Address, u64>>")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "MultiResultVec", TypeParameters: []*TypeFormula{
		{Name: "MultiResult2", TypeParameters: []*TypeFormula{
			{Name: "Address", TypeParameters: []*TypeFormula{}},
			{Name: "u64", TypeParameters: []*TypeFormula{}},
		}},
	}}, *formula)
	require.Equal(t, "MultiResultVec<MultiResult2<Address, u64>>", formula.String())

	formula, err = parser.ParseExpression("tuple3<i32, bytes, Option<i64>>")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "tuple3", TypeParameters: []*TypeFormula{
		{Name: "i32", TypeParameters: []*TypeFormula{}},
		{Name: "bytes", TypeParameters: []*TypeFormula{}},
		{Name: "Option", TypeParameters: []*TypeFormula{
			{Name: "i64", TypeParameters: []*TypeFormula{}},
		}},
	}}, *formula)
	require.Equal(t, "tuple3<i32, bytes, Option<i64>>", formula.String())

	formula, err = parser.ParseExpression("tuple2<i32, i32>")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "tuple2", TypeParameters: []*TypeFormula{
		{Name: "i32", TypeParameters: []*TypeFormula{}},
		{Name: "i32", TypeParameters: []*TypeFormula{}},
	}}, *formula)
	require.Equal(t, "tuple2<i32, i32>", formula.String())

	formula, err = parser.ParseExpression("tuple2<i32,i32>  ")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "tuple2", TypeParameters: []*TypeFormula{
		{Name: "i32", TypeParameters: []*TypeFormula{}},
		{Name: "i32", TypeParameters: []*TypeFormula{}},
	}}, *formula)
	require.Equal(t, "tuple2<i32, i32>", formula.String())

	formula, err = parser.ParseExpression("tuple<List<u64>, List<u64>>")
	require.NoError(t, err)
	require.Equal(t, TypeFormula{Name: "tuple", TypeParameters: []*TypeFormula{
		{Name: "List", TypeParameters: []*TypeFormula{
			{Name: "u64", TypeParameters: []*TypeFormula{}},
		}},
		{Name: "List", TypeParameters: []*TypeFormula{
			{Name: "u64", TypeParameters: []*TypeFormula{}},
		}},
	}}, *formula)
	require.Equal(t, "tuple<List<u64>, List<u64>>", formula.String())
}

func TestTypeFormulaParser_TokenizeExpression(t *testing.T) {
	parser := NewTypeFormulaParser()

	tokens := parser.tokenizeExpression("i64")
	require.Equal(t, []string{"i64"}, tokens)

	tokens = parser.tokenizeExpression("tuple2<i32, i32>")
	require.Equal(t, []string{"tuple2", "<", "i32", ",", "i32", ">"}, tokens)
}
