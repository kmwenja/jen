package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.kmwenja.co.ke/jen"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	app := &cli.App{
		Name:  "jen",
		Usage: "Play with templates and data",
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "generate",
				Aliases: []string{"gen", "g"},
				Usage:   "apply a golang template to json data",
				Action:  genCmd,
			},
			&cli.Command{
				Name:    "markdown",
				Aliases: []string{"md", "m"},
				Usage:   "produce json from a markdown file with possible front matter",
				Action:  mdCmd,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}

func genCmd(c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("not enough arguments passed. Usage: jen gen <template> <data>")
	}

	templateFile := c.Args().Get(0)
	dataFile := c.Args().Get(1)

	var err error

	var dr io.Reader
	if dataFile == "-" {
		dr = io.Reader(os.Stdin)
	} else {
		f, err := os.Open(dataFile)
		if err != nil {
			return fmt.Errorf("could not read data from file %q: %w", dataFile, err)
		}
		defer f.Close()

		dr = io.Reader(f)
	}

	d, err := parseJSON(dr)
	if err != nil {
		return fmt.Errorf("could not parse json data: %w", err)
	}

	var t io.Reader
	f, err := os.Open(templateFile)
	if err != nil {
		return fmt.Errorf("could not open template file %q: %w", templateFile, err)
	}
	defer f.Close()

	t = io.Reader(f)

	err = jen.Gen(t, d, os.Stdout)
	if err != nil {
		return fmt.Errorf("could not generate output: %w", err)
	}

	return nil
}

func parseJSON(r io.Reader) (interface{}, error) {
	dec := json.NewDecoder(r)
	var d interface{}
	err := dec.Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("could not parse json: %w", err)
	}

	return d, nil
}

func mdCmd(c *cli.Context) error {
	var r io.Reader
	if c.NArg() == 0 {
		r = os.Stdin
	} else {
		filename := c.Args().Get(0)
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("could not open filename %q: %w", filename, err)
		}
		defer f.Close()

		r = io.Reader(f)
	}

	d, err := jen.YamlMarkdown(r)
	if err != nil {
		return fmt.Errorf("could not parse contents: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(&d)
	if err != nil {
		return fmt.Errorf("could not write json to stdout: %w", err)
	}

	return nil
}

func parseYAML(r io.Reader) (interface{}, error) {
	dec := yaml.NewDecoder(r)
	var d interface{}
	err := dec.Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("could not parse yaml: %w", err)
	}

	return d, nil
}
