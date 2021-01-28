package svg

import (
	"strings"
)

func element(name string, attrs map[string]string) *Element {
	return &Element{
		Name:       name,
		Attributes: attrs,
		Children:   []*Element{},
	}
}

func parse(svg string, validate bool) (*Element, error) {
	element, err := Parse(strings.NewReader(svg), validate)
	return element, err
}
