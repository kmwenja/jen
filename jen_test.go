package jen

import (
	"bytes"
	"testing"
)

func TestJen(t *testing.T) {
	res, err := Jen([]byte("# Hello World\nHi!"), nil)
	if err != nil {
		t.Fatalf("Expected no error, got `%v`", err)
	}
	expected := []byte("<h1>Hello World</h1>\n\n<p>Hi!</p>\n")
	if bytes.Compare(res, expected) != 0 {
		t.Fatalf("Expected %q, got %q", expected, res)
	}
}

func TestJenWithFrontmatter(t *testing.T) {
	testMd := []byte(`
	---
	title: Something
	---
	# {{ .front_matter.title }}
	This is markdown content.
    `)
	testMd = bytes.TrimSpace(testMd)
	testMd = bytes.ReplaceAll(testMd, []byte("\t"), []byte(""))
	// t.Fatalf("%q", testMd)
	res, err := Jen(testMd, nil)
	if err != nil {
		t.Fatalf("Expected no error, got `%v`", err)
	}
	expected := []byte("<h1>Something</h1>\n\n<p>This is markdown content.</p>\n")
	if bytes.Compare(res, expected) != 0 {
		t.Fatalf("Expected %q, got %q", expected, res)
	}
}

func TestJenWithContext(t *testing.T) {
	testMd := []byte(`
	# {{ .extra.title }}
	This is markdown content.
    `)
	testMd = bytes.TrimSpace(testMd)
	testMd = bytes.ReplaceAll(testMd, []byte("\t"), []byte(""))
	// t.Fatalf("%q", testMd)
	res, err := Jen(testMd, map[string]interface{}{
		"title": "Something",
	})
	if err != nil {
		t.Fatalf("Expected no error, got `%v`", err)
	}
	expected := []byte("<h1>Something</h1>\n\n<p>This is markdown content.</p>\n")
	if bytes.Compare(res, expected) != 0 {
		t.Fatalf("Expected %q, got %q", expected, res)
	}
}
