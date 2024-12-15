package markdown

import (
	"bytes"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func Convert(r io.Reader, w io.Writer) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	if err := md.Convert(buf.Bytes(), w); err != nil {
		return err
	}

	return nil
}
