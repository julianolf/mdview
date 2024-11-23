# Markdown Viewer

Open Markdown files using the default web browser for preview.

The sole purpose of this program is to provide an idea of how Markdown files will display when viewed on sites like GitHub.

## Supported Formats

- Original Markdown
- GitHub Flavored Markdown

## Installation

**Requirements:** [Go](https://go.dev) 1.16+.

To install, run:

```sh
go install github.com/julianolf/mdview@latest
```

## Usage

Preview a single file:

```sh
mdview file.md
```

You can open multiple files for preview:

```sh
mdview file1.md file2.md
```

Or use wildcards:

```sh
mdview *.md
```

For more information, run:

```sh
mdview -h
```

## Credits

- Markdown parser: [goldmark](https://github.com/yuin/goldmark)
- CSS: [github-markdown-css](https://github.com/sindresorhus/github-markdown-css)
