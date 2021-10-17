package trimers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrimers(t *testing.T) {
	p := &Predictor{}
	p.Train("class1", "one:two a ab")
	p.Train("class2", "One three")
	p.Train("class3", "one three Ещё")
	require.Equal(t, []map[rune]map[interface{}]struct{}{
		{
			156130: {
				"class1": struct{}{},
				"class2": struct{}{},
				"class3": struct{}{},
			},
			162766: {
				"class2": struct{}{},
				"class3": struct{}{},
			},
			163318: {
				"class1": struct{}{},
			},
			1516079: {
				"class3": struct{}{},
			},
		},
		{
			146695: {
				"class2": struct{}{},
				"class3": struct{}{},
			},
		},
		{
			159904: {
				"class2": struct{}{},
				"class3": struct{}{},
			},
		},
	}, p.p.m, "Check model")
	assert.Equal(t, []interface{}{"class3", "class2"}, p.Predict("tree еще"))
}

