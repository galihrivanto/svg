package svg

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
)

var sortAttributes = false

// SetSortAttributes override value whether we need
// to sort attribute during serialization
func SetSortAttributes(v bool) {
	sortAttributes = v
}

// Serialize serializes element
func (e *Element) Serialize() xml.StartElement {
	var attributes []xml.Attr

	if sortAttributes {
		// get keys
		var keys []string
		for key := range e.Attributes {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			value := e.Attributes[key]
			attr := xml.Attr{
				Name:  xml.Name{Local: key},
				Value: value,
			}
			attributes = append(attributes, attr)
		}

	} else {
		for name, value := range e.Attributes {
			attr := xml.Attr{
				Name:  xml.Name{Local: name},
				Value: value,
			}
			attributes = append(attributes, attr)
		}
	}

	return xml.StartElement{
		Name: xml.Name{Local: e.Name},
		Attr: attributes,
	}
}

// Encode encodes the element
func (e *Element) Encode(encoder *xml.Encoder) error {
	start := e.Serialize()

	if err := encoder.EncodeToken(start); err != nil {
		return err
	}
	end := start.End()

	var content xml.Token

	content = xml.CharData(e.Content)
	encoder.EncodeToken(content)

	for _, child := range e.Children {
		if err := child.Encode(encoder); err != nil {
			return err
		}
	}
	return encoder.EncodeToken(end)
}

// Render renders element to SVG
func Render(e *Element, w io.Writer) error {
	encoder := xml.NewEncoder(w)

	if err := e.Encode(encoder); err != nil {
		return fmt.Errorf("Could not render element: %s", err)
	}

	return encoder.Flush()
}
