package svg

import (
	"log"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	var testCases = []struct {
		svg     string
		element Element
	}{
		{
			`
		<svg width="100" height="100">
			<circle cx="50" cy="50" r="40" fill="red" />
		</svg>
		`,
			Element{
				Name: "svg",
				Attributes: map[string]string{
					"width":  "100",
					"height": "100",
				},
				Children: []*Element{
					element("circle", map[string]string{"cx": "50", "cy": "50", "r": "40", "fill": "red"}),
				},
			},
		},
		{
			`
		<svg height="400" width="450">
			<g stroke="black" stroke-width="3" fill="black">
				<path id="AB" d="M 100 350 L 150 -300" stroke="red" />
				<path id="BC" d="M 250 50 L 150 300" stroke="red" />
				<path d="M 175 200 L 150 0" stroke="green" />
			</g>
		</svg>
		`,
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
							"stroke":       "black",
							"stroke-width": "3",
							"fill":         "black",
						},
						Children: []*Element{
							element("path", map[string]string{"id": "AB", "d": "M 100 350 L 150 -300", "stroke": "red"}),
							element("path", map[string]string{"id": "BC", "d": "M 250 50 L 150 300", "stroke": "red"}),
							element("path", map[string]string{"d": "M 175 200 L 150 0", "stroke": "green"}),
						},
					},
				},
			},
		},
		{
			"",
			Element{},
		},
	}

	for _, test := range testCases {
		actual, err := parse(test.svg, false)

		if !(test.element.Compare(actual) && err == nil) {
			t.Errorf("Parse: expected %v, actual %v\n", test.element, actual)
		}
	}
}

func TestValidDocument(t *testing.T) {
	svg := `
		<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="svg-root" width="100%" height="100%" viewBox="0 0 480 360">
			<title id="test-title">color-prop-01-b</title>
			<desc id="test-desc">Test that viewer has the basic capability to process the color property</desc>
			<rect id="test-frame" x="1" y="1" width="478" height="358" fill="none" stroke="#000000"/>
		</svg>
		`

	element, err := parse(svg, true)
	if !(element != nil && err == nil) {
		t.Errorf("Validation: expected %v, actual %v\n", nil, err)
	}
}

func TestParseInkscapeSVG(t *testing.T) {
	f, err := os.Open("./inkscape.svg")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer f.Close()

	root, err := Parse(f, false)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	log.Println(root.Attributes)
}
