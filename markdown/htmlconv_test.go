package markdown

import (
	"bytes"
	"testing"
)

const md string = `# Head 1

Some text.

- Item 1
- Item 2
  - Sub Item 1
  - Sub Item 2
- Item 3
`

const html = `<h1>Head 1</h1>
<p>Some text.</p>
<ul>
<li>Item 1</li>
<li>Item 2
<ul>
<li>Sub Item 1</li>
<li>Sub Item 2</li>
</ul>
</li>
<li>Item 3</li>
</ul>
`

func TestConvert(t *testing.T) {
	r := bytes.NewReader([]byte(md))
	var w bytes.Buffer

	if err := Convert(r, &w); err != nil {
		t.Fatalf("Fail to convert: %v\n", err)
	}

	got := w.String()
	if got != html {
		t.Fatalf("Conversion mismatch: %v != %v\n", got, html)
	}
}
