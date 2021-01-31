package svg

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// define library errors
var (
	ErrElementNotFound = errors.New("element not found")
)

// SetContent replace element text by id
func SetContent(root *Element, id string, text string) error {
	el := root.FindID(id)
	if el == nil {
		return ErrElementNotFound
	}
	el.Content = text

	return nil
}

// Set64Image replace element embeded image by id
func Set64Image(root *Element, id string, content string) error {
	el := root.FindID(id)
	if el == nil {
		return ErrElementNotFound
	}

	// check whether content is already valid base64
	// image format
	if !strings.HasPrefix(content, "data:") {
		mime := http.DetectContentType([]byte(content))
		content = fmt.Sprintf("data:%s;base64,%s", mime, content)
	}

	el.Attributes["href"] = content

	return nil
}

// SetImage replace image by id
func SetImage(root *Element, id string, path string, embed bool) error {
	// open image
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	b, err := ioutil.ReadAll(fi)
	if err != nil {
		return err
	}

	imageContent := base64.StdEncoding.EncodeToString(b)

	return Set64Image(root, id, imageContent)
}
