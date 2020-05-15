package jen

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
)

// Jen generates html from `markdown` and injects it into an optional template
// specified in its yaml front matter. The rest of the front matter can then
// be used as template context. Extra template context can be added with `context` map.
func Jen(markdown []byte, context map[string]interface{}) ([]byte, error) {
	data := make(map[string]interface{})
	templateFiles := make([]string, 0)

	if len(context) > 0 {
		data["extra"] = context
	}

	// TODO support json and toml frontmatter
	yml, md := ParseFrontmatter(markdown)
	if len(yml) > 0 {
		var yamlData map[string]interface{}
		err := yaml.Unmarshal(yml, &yamlData)
		if err != nil {
			return nil, fmt.Errorf("could not parse yaml frontmatter: %v", err)
		}
		data["front_matter"] = yamlData
		if val, ok := yamlData["template"]; ok {
			switch val.(type) {
			case string:
				templateFiles = append(templateFiles, val.(string))
			case []interface{}:
				for _, i := range val.([]interface{}) {
					s, ok := i.(string)
					if !ok {
						return nil, fmt.Errorf("unsupported template value: (%T) %v", val, val)
					}
					templateFiles = append(templateFiles, s)
				}
			default:
				return nil, fmt.Errorf("unsupported template value: (%T) %v", val, val)
			}
		}
	}

	html := blackfriday.Run(md)

	localTpl, err := template.New("__default__").Parse(string(html))
	if err != nil {
		return nil, fmt.Errorf("could not parse template in markdown: %v", err)
	}

	buf := &bytes.Buffer{}
	err = localTpl.Execute(buf, data)
	if err != nil {
		return nil, fmt.Errorf("could not apply template in markdown: %v", err)
	}

	if len(templateFiles) == 0 {
		return removeNoValue(buf.Bytes()), nil
	}
	tpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		return nil, fmt.Errorf("could not parse template files `%v`: %v", templateFiles, err)
	}

	data["content"] = buf.String()
	buf.Truncate(0)
	if tpl.Lookup("template") != nil {
		err = tpl.ExecuteTemplate(buf, "template", data)
	} else {
		err = tpl.Execute(buf, data)
	}
	if err != nil {
		return nil, fmt.Errorf("could not apply template: %v", err)
	}

	return removeNoValue(buf.Bytes()), nil
}

func removeNoValue(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("<no value>"), []byte(""))
}
