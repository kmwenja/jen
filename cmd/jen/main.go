package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"git.kmwenja.co.ke/jen"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	app := &cli.App{
		Name:  "jen",
		Usage: "Generate a html page from markdown",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Value:   nil,
				Usage:   "extra data to use as template context",
			},
		},
		Action: func(c *cli.Context) error {
			var data []byte
			var err error

			if c.NArg() == 0 {
				data, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("could not read from stdin: %v", err)
				}
			} else {
				path := c.Args().Get(0)
				data, err = ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("could not read from path `%s`: %v", path, err)
				}
			}

			context := make(map[string]interface{})
			dataPaths := c.StringSlice("data")
			if len(dataPaths) > 0 {
				// TODO support toml files
				for _, p := range dataPaths {
					var d map[string]interface{}
					var err error
					ext := path.Ext(p)
					switch ext {
					case ".json":
						// TODO make an internal utility package for this parsing
						// and for the map merging
						d, err = parseJSONFile(p)
					case ".yaml":
						d, err = parseYAMLFile(p)
					default:
						return fmt.Errorf("unsupported data filetype %q: %s", ext, p)
					}
					if err != nil {
						return fmt.Errorf("could not parse %s file %q: %v", ext, p, err)
					}
					for k, v := range d {
						context[k] = v
					}
				}
			}

			html, err := jen.Jen(data, context)
			if err != nil {
				return fmt.Errorf("could not generate html: %v", err)
			}

			_, err = fmt.Fprintf(os.Stdout, "%s", html)
			if err != nil {
				return fmt.Errorf("could not write to stdout: %v", err)
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if err != nil {
			panic(err)
		}
	}
}

func parseJSONFile(filepath string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %q: %v", filepath, err)
	}

	var d map[string]interface{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, fmt.Errorf("could not parse json in %q: %v", filepath, err)
	}

	return d, nil
}

func parseYAMLFile(filepath string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %q: %v", filepath, err)
	}

	var d map[string]interface{}
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		return nil, fmt.Errorf("could not parse json in %q: %v", filepath, err)
	}

	return d, nil
}
