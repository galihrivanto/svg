package utils

import (
	"testing"
)

func TestStyleParser(t *testing.T) {
	var testCases = []struct {
		styles   string
		expected Styles
	}{
		{
			"fill:white;stroke:#000000;",
			Styles(
				[]*Style{
					{Property: "fill", Value: "white"},
					{Property: "stroke", Value: "#000000"},
				},
			),
		},
		{
			"fill:white;stroke-opacity:1",
			Styles(
				[]*Style{
					{Property: "fill", Value: "white"},
					{Property: "stroke-opacity", Value: "1"},
				},
			),
		},
	}

	for _, test := range testCases {
		styles := StyleParser(test.styles)
		if !test.expected.Compare(styles) {
			t.Errorf("Style: expected %v, actual %v\n", test.expected, styles)
		}
	}
}
