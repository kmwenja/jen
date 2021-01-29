package jen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
)

// Gen applies a go text/template.Template parsed from `templateContents` to the
// `data` provided and writes that to `output`.
func Gen(templateFiles []string, data interface{}, output io.Writer) error {
	t := template.New("__default__")
	for _, filename := range templateFiles {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("could not read template file %q: %w", filename, err)
		}
		_, err = t.Parse(string(contents))
		if err != nil {
			return fmt.Errorf("could not parse template file %q: %w", filename, err)
		}
	}

	err := t.Execute(output, data)
	if err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	return nil
}

// YamlMarkdown takes yaml-markdown contents from `r` and returns a
// map[string]interface that contains the yaml inside the contents as well
// as parsed markdown from the contents as a map of strings to interface{}.
// The yaml will be under the key `yaml` while the markdown will be under the key `markdown`
func YamlMarkdown(r io.Reader) (map[string]interface{}, error) {
	yml, md, err := ParseFrontmatter(r)
	if err != nil {
		return nil, fmt.Errorf("could not parse contents: %w", err)
	}

	data := make(map[string]interface{})

	if len(yml) > 0 {
		var yamlData map[string]interface{}
		err := yaml.Unmarshal(yml, &yamlData)
		if err != nil {
			return nil, fmt.Errorf("could not parse yaml frontmatter: %w", err)
		}
		data["yaml"] = yamlData
	}

	data["markdown"] = string(blackfriday.Run(md))
	return data, nil
}

// ParseFrontmatter extracts frontmatter from `r` (representing a
// reader with yaml-markdown) and returns the frontmatter and the
// markdown
func ParseFrontmatter(r io.Reader) (matter []byte, rest []byte, err error) {
	var (
		beforematter = 1
		inmatter     = 2
		aftermatter  = 3
	)
	status := beforematter

	hyphenCount := 0

	mainBuf := bufio.NewReader(r)
	matterBuf := &bytes.Buffer{}
	mdBuf := &bytes.Buffer{}
	tempBuf := &bytes.Buffer{}

	for {
		b, err := mainBuf.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, fmt.Errorf("could not read byte: %w", err)
		}

		switch status {
		case beforematter:
			switch b {
			case '-':
				hyphenCount++
				// just in case it turns out it wasn't a separator
				tempBuf.WriteByte(b)
			case '\n':
				if hyphenCount >= 3 {
					status = inmatter
					hyphenCount = 0
					// since we got front matter, reset the content
					// we had already put into the markdown buffer
					mdBuf.Truncate(0)
					tempBuf.Truncate(0)
				} else {
					mdBuf.Write(tempBuf.Bytes())
					tempBuf.Truncate(0)
					mdBuf.WriteByte(b)
				}
			default:
				// assume there's no front matter
				mdBuf.Write(tempBuf.Bytes())
				tempBuf.Truncate(0)
				mdBuf.WriteByte(b)
				hyphenCount = 0
			}
		case inmatter:
			switch b {
			case '-':
				hyphenCount++
				// just in case it turns out it wasn't a separator
				tempBuf.WriteByte(b)
			case '\n':
				if hyphenCount >= 3 {
					status = aftermatter
					hyphenCount = 0
					tempBuf.Truncate(0)
				} else {
					matterBuf.Write(tempBuf.Bytes())
					tempBuf.Truncate(0)
					matterBuf.WriteByte(b)
				}
			default:
				matterBuf.Write(tempBuf.Bytes())
				tempBuf.Truncate(0)
				matterBuf.WriteByte(b)
				hyphenCount = 0
			}
		case aftermatter:
			mdBuf.WriteByte(b)
		}
	}

	// in case any temp stuff was left
	switch status {
	case beforematter, aftermatter:
		mdBuf.Write(tempBuf.Bytes())
	case inmatter:
		matterBuf.Write(tempBuf.Bytes())
	}

	matter = matterBuf.Bytes()
	rest = mdBuf.Bytes()
	err = nil
	return
}
