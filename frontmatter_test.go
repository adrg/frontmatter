package frontmatter_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

func TestMatters(t *testing.T) {
	// Test structs.
	type (
		metadata struct {
			Size int `yaml:"size" toml:"size" json:"size"`
		}

		matter struct {
			Name     string   `yaml:"name" toml:"name" json:"name"`
			Tags     []string `yaml:"tags" toml:"tags" json:"tags"`
			Metadata metadata `yaml:"metadata" toml:"metadata" json:"metadata"`
		}

		parseFunc func(r io.Reader, v interface{},
			formats ...*frontmatter.Format) ([]byte, error)
		expFunc func(in string) (mExp *matter, rExp string, eExp bool)

		testCase struct {
			input        string
			formats      []*frontmatter.Format
			expParse     expFunc
			expMustParse expFunc
		}
	)

	// Expected functions.
	var (
		expValidMatter = func(input string) (*matter, string, bool) {
			return &matter{
				Name: "frontmatter",
				Tags: []string{"go", "yaml", "json", "toml"},
				Metadata: metadata{
					Size: 10,
				},
			}, "rest of the file", false
		}
		expNoRest = func(input string) (*matter, string, bool) {
			mExp, _, eExp := expValidMatter(input)
			return mExp, "", eExp
		}
		expEmptyMatter = func(input string) (*matter, string, bool) {
			return &matter{}, "", false
		}
		expNoMatter = func(input string) (*matter, string, bool) {
			return &matter{}, input, false
		}
		expMatterErr = func(input string) (*matter, string, bool) {
			return &matter{}, "", true
		}
	)

	// Test data.
	testCases := []*testCase{
		// -----------------
		// - Valid matter. -
		// -----------------

		{
			input: `
---
name: "frontmatter"
tags:
  - "go"
  - "yaml"
  - "json"
  - "toml"
metadata:
  size: 10
---
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
---yaml
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
metadata:
  size: 10
---
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
+++
name = "frontmatter"
tags = ["go", "yaml", "json", "toml"]
[metadata]
size = 10
+++
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
---toml
name = "frontmatter"
tags = ["go", "yaml", "json", "toml"]
[metadata]
size = 10
---
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
;;;
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
	"size": 10
  }
}
;;;
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
---json
{
  "name": "frontmatter",
  "metadata": {
	"size": 10
  },
  "tags": ["go", "yaml", "json", "toml"]
}
---
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
{
  "name": "frontmatter",
  "metadata": {
	"size": 10
  },
  "tags": ["go", "yaml", "json", "toml"]
}

rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
{
  "metadata": {
	"size": 10
  },
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"]
}

rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
	"size": 10
  }
}

rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
		},

		{
			input: `
---
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
metadata:
  size: 10
---
`,
			expParse:     expNoRest,
			expMustParse: expNoRest,
		},

		{
			input: `
---yaml
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
metadata:
  size: 10
---`,
			expParse:     expNoRest,
			expMustParse: expNoRest,
		},

		{
			input: `{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
    "size": 10
  }
}
`,
			expParse:     expNoRest,
			expMustParse: expNoRest,
		},

		{
			input: `
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
    "size": 10
  }
}`,
			expParse:     expNoRest,
			expMustParse: expNoRest,
		},

		// -----------------
		// - Empty matter. -
		// -----------------

		{
			input: `---
	 ---`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `---yaml
	 ---`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `+++
	 +++`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `---toml
	 ---`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `;;;
			{}
	 ;;;`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `---json
			{}
	 ---`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `{
	 }

`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `
{
	 }
`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `{
	 }`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		{
			input: `{
	 }`,
			expParse:     expEmptyMatter,
			expMustParse: expEmptyMatter,
		},

		// --------------
		// - No matter. -
		// --------------

		{
			input:        ``,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input:        `rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
	`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `start of file
---
name: "frontmatter"
tags:
  - "go"
  - "yaml"
  - "json"
  - "toml"
metadata:
  size: 10
---
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `start of file
---yaml
name: "frontmatter"
tags:
  - "go"
  - "yaml"
  - "json"
  - "toml"
metadata:
  size: 10
---
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
start of file
+++
name = "frontmatter"
tags = ["go", "yaml", "json", "toml"]
[metadata]
size = 10
+++
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
start of file
---toml
name = "frontmatter"
tags = ["go", "yaml", "json", "toml"]
[metadata]
size = 10
---
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
			start of file
;;;
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
	"size": 10
  }
}
;;;
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
			start of file
---json
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
	"size": 10
  }
}
---
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
			start of file
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"],
  "metadata": {
	"size": 10
  }
}

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		// -------------------
		// - Invalid matter. -
		// -------------------

		{
			input: `
---
name: "frontmatter"

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
---json
name: "frontmatter"

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
+++
name = "frontmatter"

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
---toml
name = "frontmatter"

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
;;;
{
  "name": "frontmatter",

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
---json
{
  "name": "frontmatter",

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
{
  "name": "frontmatter",

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
name: "frontmatter"
---
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
name = "frontmatter"
+++

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
  "name": "frontmatter",
}
;;;
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
  "name": "frontmatter",
}

rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
{
  "name": "frontmatter"
}
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		{
			input: `
{
  "name": "frontmatter"
}
should be an empty line
rest of the file`,
			expParse:     expNoMatter,
			expMustParse: expMatterErr,
		},

		// -----------------
		// - Invalid data. -
		// -----------------

		{
			input: `---
name = "frontmatter"
---
rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `---yaml
name = "frontmatter"
---
rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `+++
name: "frontmatter"
+++
rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `---toml
name: "frontmatter"
---
rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `;;;
{
  name: "frontmatter"
}
;;;
rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `---json
{
  name: "frontmatter"
}
---
rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `{
name: "frontmatter"
}

rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		{
			input: `{
"name": "frontmatter",
}

rest of the file`,
			expParse:     expMatterErr,
			expMustParse: expMatterErr,
		},

		// ------------------
		// - Custom formats -
		// ------------------

		{
			input: `
...
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
metadata:
  size: 10
...
rest of the file`,
			expParse:     expValidMatter,
			expMustParse: expValidMatter,
			formats: []*frontmatter.Format{
				frontmatter.NewFormat("...", "...", yaml.Unmarshal),
			},
		},
	}

	failFunc := func(in string, exp, act interface{}) {
		t.Fatalf("Input: `%s`\n\nNot equal: \n"+
			"expected: %v\n"+
			"actual  : %v", in, exp, act)
	}

	checkFunc := func(in string, mExp, mAct *matter,
		rExp, rAct string, eExp, eAct bool) {
		if eExp != eAct {
			failFunc(in, eExp, eAct)
		}
		if mExp.Name != mAct.Name {
			failFunc(in, mExp.Name, mAct.Name)
		}
		if strings.Join(mExp.Tags, ",") != strings.Join(mAct.Tags, ",") {
			failFunc(in, mExp.Tags, mAct.Tags)
		}
		if mExp.Metadata.Size != mAct.Metadata.Size {
			failFunc(in, mExp.Metadata.Size, mAct.Metadata.Size)
		}
		if rExp != rAct {
			failFunc(in, rExp, rAct)
		}
	}

	testFunc := func(in string, formats []*frontmatter.Format,
		expFunc expFunc, parseFunc parseFunc) {
		// Get expected data.
		mExp, rExp, hasErr := expFunc(in)

		// Get actual data.
		mAct := &matter{}
		rest, err := parseFunc(strings.NewReader(in), mAct, formats...)
		checkFunc(in, mExp, mAct, rExp, string(rest), hasErr, err != nil)
	}

	for _, tc := range testCases {
		testFunc(tc.input, tc.formats, tc.expParse, frontmatter.Parse)
		testFunc(tc.input, tc.formats, tc.expMustParse, frontmatter.MustParse)
	}
}

func TestRequiredNewLineDelimiters(t *testing.T) {
	testCases := []struct {
		name string
		text string
	}{
		{name: "plain", text: "value\n"},
		{name: "one literal delimiter", text: "...\n"},
		{name: "consecutive literal delimiters", text: "...\n...\n"},
		{name: "delimiter followed by content", text: "...\nvalue\n...\n"},
	}
	parsers := []struct {
		name  string
		parse func(io.Reader, interface{}, ...*frontmatter.Format) ([]byte, error)
	}{
		{name: "Parse", parse: frontmatter.Parse},
		{name: "MustParse", parse: frontmatter.MustParse},
	}
	for _, parser := range parsers {
		t.Run(parser.name, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					for _, lineEnding := range []string{"\n", "\r\n"} {
						t.Run(fmt.Sprintf("line ending %q", lineEnding), func(t *testing.T) {
							input := "---yaml\ntext: |\n  " + strings.ReplaceAll(
								strings.TrimSuffix(tc.text, "\n"), "\n", "\n  ") + "\n...\n\nbody"
							input = strings.ReplaceAll(input, "\n", lineEnding)
							format := frontmatter.NewFormat("---yaml", "...", yaml.Unmarshal)
							format.RequiresNewLine = true
							var matter struct {
								Text string `yaml:"text"`
							}
							rest, err := parser.parse(strings.NewReader(input), &matter, format)
							if err != nil {
								t.Fatal(err)
							}
							if matter.Text != tc.text {
								t.Errorf("expected text %q, got %q", tc.text, matter.Text)
							}
							if string(rest) != "body" {
								t.Errorf("expected body %q, got %q", "body", rest)
							}
						})
					}
				})
			}
		})
	}
}
