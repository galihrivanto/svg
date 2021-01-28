package svg

import (
	"testing"
)

func testElement() *Element {
	svg := `
		<svg width="1000" height="600">
			<g id="first">
				<rect width="5" height="3" id="inFirst"/>
				<rect width="5" height="2" id="inFirst"/>
			</g>
			<g id="second">
				<path d="M50 50 Q50 100 100 100"/>
				<rect width="5" height="1"/>
			</g>
		</svg>
	`
	element, _ := parse(svg, false)
	return element
}

func equals(t *testing.T, name string, expected, actual *Element) {
	if !(expected == actual || expected.Compare(actual)) {
		t.Errorf("%s: expected %v, actual %v\n", name, expected, actual)
	}
}

func equalSlices(t *testing.T, name string, expected, actual []*Element) {
	if len(expected) != len(actual) {
		t.Errorf("%s: expected %v, actual %v\n", name, expected, actual)
		return
	}

	for i, r := range actual {
		equals(t, name, expected[i], r)
	}
}

func TestFindAll(t *testing.T) {
	svgElement := testElement()

	equalSlices(t, "Find", []*Element{
		element("rect", map[string]string{"width": "5", "height": "3", "id": "inFirst"}),
		element("rect", map[string]string{"width": "5", "height": "2", "id": "inFirst"}),
		element("rect", map[string]string{"width": "5", "height": "1"}),
	}, svgElement.FindAll("rect"))

	equalSlices(t, "Find", []*Element{}, svgElement.FindAll("circle"))
}

func TestFindID(t *testing.T) {
	svgElement := testElement()

	equals(t, "Find", &Element{
		Name:       "g",
		Attributes: map[string]string{"id": "second"},
		Children: []*Element{
			element("path", map[string]string{"d": "M50 50 Q50 100 100 100"}),
			element("rect", map[string]string{"width": "5", "height": "1"}),
		},
	}, svgElement.FindID("second"))

	equals(t, "Find",
		element("rect", map[string]string{"width": "5", "height": "3", "id": "inFirst"}),
		svgElement.FindID("inFirst"),
	)

	equals(t, "Find", nil, svgElement.FindID("missing"))
}
