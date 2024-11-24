package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const semver = "v0.1.0"

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

func init() {
	tmp := os.TempDir()

	flag.StringVar(&oflag, "o", tmp, "Output directory for the preview file")
	flag.BoolVar(&vflag, "v", false, "Show version number and quit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [file ...]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func open(filename string) error {
	url := "file://" + filename

	// FIXME: Better handling of Linux variants, maybe support other OSs
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("Unsupported OS")
	}

	return cmd.Start()
}

func convert(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	markdown := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buffer bytes.Buffer
	if err := markdown.Convert(content, &buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
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

	if err := open(filename); err != nil {
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
	flag.Parse()
	if vflag {
		version()
	}

	run(flag.Args())
}
