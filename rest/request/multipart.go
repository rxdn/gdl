package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
)

type Attachment struct {
	Id          int    `json:"id"`
	Description string `json:"description,omitempty"`
	FileName    string `json:"filename"`
	File        File   `json:"-"`
}

type File struct {
	ContentType string
	Reader      io.Reader
}

type MultipartPayload interface {
	GetAttachments() []Attachment
}

func EncodeMultipartFormData(payload MultipartPayload) ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	if err := writer.WriteField("payload_json", string(payloadJson)); err != nil {
		return nil, "", err
	}

	for i, file := range payload.GetAttachments() {
		fileName := file.FileName
		fileName = strings.Replace(fileName, "\\", "\\\\", -1)
		fileName = strings.Replace(fileName, "\"", "\\\"", -1)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files[%d]"; filename="%s"`, i, fileName))
		h.Set("Content-Type", file.File.ContentType)

		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, "", err
		}

		if _, err := io.Copy(part, file.File.Reader); err != nil {
			return nil, "", err
		}
	}

	return []byte(string(body.Bytes()) + "\r\n--" + writer.Boundary() + "--"), writer.Boundary(), nil
}
