package utils

import (
	"testing"
)

func TestPathParser(t *testing.T) {
	var testCases = []struct {
		d        string
		expected *Path
	}{
		{
			"M 10,20 L 30,30 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{10, 20}},
							{Symbol: "L", Params: []float64{30, 30}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
		{
			"M .2.3 L 30,30 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{0.2, 0.3}},
							{Symbol: "L", Params: []float64{30, 30}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
		{
			"M10-20 L30,30 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{10, -20}},
							{Symbol: "L", Params: []float64{30, 30}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
		{
			"M 10-20 L 30,30 L 40,40 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{10, -20}},
							{Symbol: "L", Params: []float64{30, 30}},
							{Symbol: "L", Params: []float64{40, 40}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
		{
			"M10,20 L20,30 L10,20",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{10, 20}},
							{Symbol: "L", Params: []float64{20, 30}},
							{Symbol: "L", Params: []float64{10, 20}},
						},
					},
				},
			},
		},
	}
	for _, test := range testCases {
		path, err := PathParser(test.d)
		if !(test.expected.Compare(path) && err == nil) {
			t.Errorf("Path: expected %v, actual %v\n", test.expected, path)
		}
	}
}

func TestParamNumberInPath(t *testing.T) {
	path, err := PathParser("M 10 20 30 Z")
	expectedError := "Incorrect number of parameters for M"

	if !(path == nil && err.Error() == expectedError) {
		t.Errorf("Path: expected %v, actual %v\n", expectedError, err)
	}
}

func TestMissingZero(t *testing.T) {
	var testCases = []struct {
		d        string
		expected *Path
	}{
		{
			"M 0.2 0.3 L 30,30 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{0.2, 0.3}},
							{Symbol: "L", Params: []float64{30, 30}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		path, err := PathParser(test.d)
		if !(test.expected.Compare(path) && err == nil) {
			t.Errorf("Path: expected %v, actual %v\n", test.expected, path)
		}
	}
}

func TestTwoSubpaths(t *testing.T) {
	var testCases = []struct {
		d        string
		expected *Path
	}{
		{
			"M25,0 L0,30 L50,50 Z m 10, 10 L50,50 l10,0 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{25, 0}},
							{Symbol: "L", Params: []float64{0, 30}},
							{Symbol: "L", Params: []float64{50, 50}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
					{
						Commands: []*Command{
							{Symbol: "m", Params: []float64{10, 10}},
							{Symbol: "L", Params: []float64{50, 50}},
							{Symbol: "l", Params: []float64{10, 0}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		path, err := PathParser(test.d)
		if len(test.expected.Subpaths) != len(path.Subpaths) {
			t.Errorf("Incorrect number of subpaths found")
		}

		if !(test.expected.Compare(path) && err == nil) {
			t.Errorf("Path: expected %v, actual %v\n", *(test.expected), *path)
		}
	}
}

func TestImplicitLineCommands(t *testing.T) {
	var testCases = []struct {
		d        string
		expected *Path
	}{
		{
			"M 10,20 30,40 Z m 10,20 30,40 Z",
			&Path{
				Subpaths: []*Subpath{
					{
						Commands: []*Command{
							{Symbol: "M", Params: []float64{10, 20}},
							{Symbol: "L", Params: []float64{30, 40}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
					{
						Commands: []*Command{
							{Symbol: "m", Params: []float64{10, 20}},
							{Symbol: "l", Params: []float64{30, 40}},
							{Symbol: "Z", Params: []float64{}},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		path, err := PathParser(test.d)
		if !(test.expected.Compare(path) && err == nil) {
			t.Errorf("Path: expected %v, actual %v\n", test.expected, path)
		}
	}
}
