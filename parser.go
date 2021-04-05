package svg

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html/charset"
)

// ValidationError contains errors which have occured when parsing svg input.
type ValidationError struct {
	msg string
}

func (err ValidationError) Error() string {
	return err.msg
}

// Element is a representation of an SVG element.
type Element struct {
	Name       string
	Attributes map[string]string
	Children   []*Element
	Content    string

	// namespace dict extracted form root element attrributes
	// which prefixed with xmlns:
	Namespaces map[string]string
}

func canonizedName(name xml.Name, namespaces map[string]string) string {
	if namespaces == nil {
		return name.Local
	}

	// lookup tag key
	tag := name.Local
	if name.Space != "" {
		prefix := namespaces[name.Space]

		// ignore default svg prefix
		if prefix != "svg" {
			tag = prefix + ":" + name.Local
		}
	}

	return tag
}

// NewElement creates element from decoder token.
func NewElement(root *Element, token xml.StartElement) *Element {
	element := &Element{}
	attributes := make(map[string]string)
	var namespaces map[string]string
	if root == nil {
		namespaces = make(map[string]string)
	} else {
		namespaces = root.Namespaces
	}

	for _, attr := range token.Attr {
		key := attr.Name.Local
		if attr.Name.Space == "xmlns" {
			namespaces[attr.Value] = attr.Name.Local
			key = attr.Name.Space + ":" + attr.Name.Local
		} else if attr.Name.Space != "" {
			tag, ok := namespaces[attr.Name.Space]
			if ok {
				key = tag + ":" + attr.Name.Local
			}
		}

		attributes[key] = attr.Value
	}

	element.Name = canonizedName(token.Name, namespaces)
	element.Attributes = attributes
	if root == nil {
		element.Namespaces = namespaces
	}

	return element
}

// Compare compares two elements.
func (e *Element) Compare(o *Element) bool {
	if e.Name != o.Name || e.Content != o.Content ||
		len(e.Attributes) != len(o.Attributes) ||
		len(e.Children) != len(o.Children) {
		return false
	}

	for k, v := range e.Attributes {
		if v != o.Attributes[k] {
			return false
		}
	}

	for i, child := range e.Children {
		if !child.Compare(o.Children[i]) {
			return false
		}
	}
	return true
}

// DecodeFirst creates the first element from the decoder.
func DecodeFirst(decoder *xml.Decoder) (*Element, error) {
	for {
		token, err := decoder.Token()
		if token == nil && err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch element := token.(type) {
		case xml.StartElement:
			return NewElement(nil, element), nil
		}
	}
	return &Element{}, nil
}

// Decode decodes the child elements of element.
func (e *Element) Decode(root *Element, decoder *xml.Decoder) error {
	for {
		token, err := decoder.Token()
		if token == nil && err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch element := token.(type) {
		case xml.StartElement:
			nextElement := NewElement(root, element)
			err := nextElement.Decode(root, decoder)
			if err != nil {
				return err
			}

			e.Children = append(e.Children, nextElement)

		case xml.CharData:
			data := strings.TrimSpace(string(element))
			if data != "" {
				e.Content = string(element)
			}

		case xml.EndElement:
			namespaces := e.Namespaces
			if root != nil {
				namespaces = root.Namespaces
			}

			if canonizedName(element.Name, namespaces) == e.Name {
				return nil
			}
		}
	}
	return nil
}

// Parse creates an Element instance from an SVG input.
func Parse(source io.Reader, validate bool) (*Element, error) {
	raw, err := ioutil.ReadAll(source)
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	decoder.CharsetReader = charset.NewReaderLabel
	element, err := DecodeFirst(decoder)
	if err != nil {
		return nil, err
	}
	if err := element.Decode(element, decoder); err != nil && err != io.EOF {
		return nil, err
	}
	return element, nil
}
