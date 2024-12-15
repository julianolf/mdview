package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/browser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const semver = "v0.2.0"

//go:embed github-markdown.css
var css string

//go:embed template.html
var html string

var (
	tmpl *template.Template = template.Must(template.New("html").Parse(html))
	wg   sync.WaitGroup
)

var (
	oflag string
	vflag bool
)

type Data struct {
	CSS     template.CSS
	Content template.HTML
}

func version() {
	fmt.Println(os.Args[0], semver)
	os.Exit(0)
}

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: %s [file ...]\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.StringVar(&oflag, "o", os.TempDir(), "Output directory for the preview file")
	flag.BoolVar(&vflag, "v", false, "Show version number and quit")
	flag.Parse()
}

func convert(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var in bytes.Buffer
	if _, err := io.Copy(&in, file); err != nil {
		return "", err
	}

	var out bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	if err := md.Convert(in.Bytes(), &out); err != nil {
		return "", err
	}

	return out.String(), nil
}

func write(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := Data{template.CSS(css), template.HTML(content)}
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}

func rename(filename string) string {
	filename = filepath.Base(filename)
	dirname, _ := filepath.Abs(oflag)
	prefix, _ := strings.CutSuffix(filename, ".md")
	return filepath.Join(dirname, prefix+".html")
}

func preview(filename string) {
	defer wg.Done()

	content, err := convert(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	filename = rename(filename)
	if err := write(filename, content); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := browser.OpenFile(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(files []string) {
	for _, filename := range files {
		wg.Add(1)
		go preview(filename)
	}
	wg.Wait()
}

func main() {
	if vflag {
		version()
	}

	run(flag.Args())
}
