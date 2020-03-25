package rest

import (
	"encoding/base64"
	"fmt"
	"github.com/rxdn/gdl/rest/request"
	"io"
	"io/ioutil"
)

type Image struct {
	ContentType request.ContentType
	ImageReader io.Reader
}

func (i *Image) Encode() (string, error) {
	content, err := ioutil.ReadAll(i.ImageReader)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	return fmt.Sprintf("data:%s;base64,%s", string(i.ContentType), encoded), nil
}

func (i Image) MarshalJSON() ([]byte, error) {
	imageData, err := i.Encode(); if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("\"%s\"", imageData)), nil
}

