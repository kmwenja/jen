package jen

import (
	"bytes"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	type testCase struct {
		Input        []byte
		ExpectedYaml []byte
		ExpectedMd   []byte
	}

	cases := []testCase{
		// without front matter
		testCase{[]byte("some stuff"), nil, []byte("some stuff")},
		testCase{[]byte("some stuff\nsome other stuff"), nil, []byte("some stuff\nsome other stuff")},

		// with invalid separators
		testCase{[]byte("-"), nil, []byte("-")},
		testCase{[]byte("--"), nil, []byte("--")},
		testCase{[]byte("---"), nil, []byte("---")},
		testCase{[]byte("-b--"), nil, []byte("-b--")},
		testCase{[]byte("-some stuff--"), nil, []byte("-some stuff--")},
		testCase{[]byte("--b-"), nil, []byte("--b-")},
		testCase{[]byte("--some stuff-"), nil, []byte("--some stuff-")},
		testCase{[]byte("---b"), nil, []byte("---b")},
		testCase{[]byte("---some stuff"), nil, []byte("---some stuff")},
		testCase{[]byte("---\n"), nil, nil},
		testCase{[]byte("---some stuff\n"), nil, []byte("---some stuff\n")},

		// front matter only
		testCase{[]byte("---\n-"), []byte("-"), nil},
		testCase{[]byte("---\n--"), []byte("--"), nil},
		testCase{[]byte("---\n---"), []byte("---"), nil},
		testCase{[]byte("---\n-b--"), []byte("-b--"), nil},
		testCase{[]byte("---\n--b-"), []byte("--b-"), nil},
		testCase{[]byte("---\n---b"), []byte("---b"), nil},
		testCase{[]byte("---\n-some stuff--"), []byte("-some stuff--"), nil},
		testCase{[]byte("---\n--some stuff-"), []byte("--some stuff-"), nil},
		testCase{[]byte("---\n---some stuff"), []byte("---some stuff"), nil},
		testCase{[]byte("---\nsome stuff"), []byte("some stuff"), nil},
		testCase{[]byte("---\nsome stuff\n"), []byte("some stuff\n"), nil},
		testCase{[]byte("---\nsome stuff\nother stuff"), []byte("some stuff\nother stuff"), nil},
		testCase{[]byte("---\nsome stuff\n-"), []byte("some stuff\n-"), nil},
		testCase{[]byte("---\nsome stuff\n--"), []byte("some stuff\n--"), nil},
		testCase{[]byte("---\nsome stuff\n---"), []byte("some stuff\n---"), nil},
		testCase{[]byte("---\nsome stuff\n---\n"), []byte("some stuff\n"), nil},
		testCase{[]byte("nonsense before---\nsome stuff\n---\n"), []byte("some stuff\n"), nil},
		testCase{[]byte("nonsense before\n---\nsome stuff\n---\n"), []byte("some stuff\n"), nil},

		// with front matter and markdown
		testCase{[]byte("---\nsome stuff\n---\nsome more stuff"), []byte("some stuff\n"), []byte("some more stuff")},
	}

	for _, c := range cases {
		yamlRes, mdRes := ParseFrontmatter(c.Input)
		if bytes.Compare(yamlRes, c.ExpectedYaml) != 0 {
			t.Fatalf("For case %q, expected yaml %q, got %q", c.Input, c.ExpectedYaml, yamlRes)
		}
		if bytes.Compare(mdRes, c.ExpectedMd) != 0 {
			t.Fatalf("For case %q, expected markdown %q, got %q", c.Input, c.ExpectedMd, mdRes)
		}
	}
}
