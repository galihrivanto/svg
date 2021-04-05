package svg

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func render(e *Element) (string, error) {
	w := &bytes.Buffer{}
	if err := Render(e, w, false); err != nil {
		return "", err
	}

	return w.String(), nil
}

func TestRender(t *testing.T) {
	var testCases = []struct {
		svg     string
		element Element
	}{
		{
			`<svg height="100" width="100"><circle cx="50" cy="50" fill="red" r="40"></circle></svg>`,
			Element{
				Name: "svg",
				Attributes: map[string]string{
					"width":  "100",
					"height": "100",
				},
				Children: []*Element{
					{
						Name:       "circle",
						Attributes: map[string]string{"cx": "50", "cy": "50", "r": "40", "fill": "red"},
					},
				},
			},
		},
		{
			`<svg height="400" width="450"><g fill="black" stroke="black" stroke-width="3"><path d="M 100 350 L 150 -300" id="AB" stroke="red"></path></g></svg>`,
			Element{
				Name: "svg",
				Attributes: map[string]string{
					"width":  "450",
					"height": "400",
				},
				Children: []*Element{
					{
						Name: "g",
						Attributes: map[string]string{
							"stroke-width": "3",
							"stroke":       "black",
							"fill":         "black",
						},
						Children: []*Element{
							{
								Name:       "path",
								Attributes: map[string]string{"id": "AB", "d": "M 100 350 L 150 -300", "stroke": "red"},
							},
						},
					},
				},
			},
		},
	}

	SetSortAttributes(true)
	for _, test := range testCases {
		actual, err := render(&test.element)

		if !(strings.Compare(test.svg, actual) == 0 && err == nil) {
			t.Errorf("Render: expected %v, actual %v\n", test.svg, actual)
		}
	}
}

func TestParseAndRenderInkscapeSVG(t *testing.T) {
	fin, err := os.Open("./inkscape.svg")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fin.Close()

	root, err := Parse(fin, false)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	fout, err := os.Create("./inkscape-result.svg")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fout.Close()

	if err := Render(root, fout, true); err != nil {
		t.Log(err)
		t.FailNow()
	}
}
