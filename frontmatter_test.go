package frontmatter_test

import (
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

		testCase struct {
			input        string
			expectedFunc func(in string) (mExp *matter, rExp string, eExp bool)
			formats      []*frontmatter.Format
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expValidMatter,
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
			expectedFunc: expNoRest,
		},

		{
			input: `
---yaml
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
metadata:
  size: 10
---`,
			expectedFunc: expNoRest,
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
			expectedFunc: expNoRest,
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
			expectedFunc: expNoRest,
		},

		// -----------------
		// - Empty matter. -
		// -----------------

		{
			input: `---
	 ---`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `---yaml
	 ---`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `+++
	 +++`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `---toml
	 ---`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `;;;
			{}
	 ;;;`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `---json
			{}
	 ---`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `{
	 }

`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `
{
	 }
`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `{
	 }`,
			expectedFunc: expEmptyMatter,
		},

		{
			input: `{
	 }`,
			expectedFunc: expEmptyMatter,
		},

		// --------------
		// - No matter. -
		// --------------

		{
			input:        ``,
			expectedFunc: expNoMatter,
		},

		{
			input:        `rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
	`,
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
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
			expectedFunc: expNoMatter,
		},

		// -------------------
		// - Invalid matter. -
		// -------------------

		{
			input: `
---
name: "frontmatter"

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
---json
name: "frontmatter"

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
+++
name = "frontmatter"

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
---toml
name = "frontmatter"

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
;;;
{
  "name": "frontmatter",

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
---json
{
  "name": "frontmatter",

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
{
  "name": "frontmatter",

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
name: "frontmatter"
---
rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
name = "frontmatter"
+++

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
  "name": "frontmatter",
}
;;;
rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
  "name": "frontmatter",
}

rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
{
  "name": "frontmatter"
}
rest of the file`,
			expectedFunc: expNoMatter,
		},

		{
			input: `
{
  "name": "frontmatter"
}
should be an empty line
rest of the file`,
			expectedFunc: expNoMatter,
		},

		// -----------------
		// - Invalid data. -
		// -----------------

		{
			input: `---
name = "frontmatter"
---
rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `---yaml
name = "frontmatter"
---
rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `+++
name: "frontmatter"
+++
rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `---toml
name: "frontmatter"
---
rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `;;;
{
  name: "frontmatter"
}
;;;
rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `---json
{
  name: "frontmatter"
}
---
rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `{
name: "frontmatter"
}

rest of the file`,
			expectedFunc: expMatterErr,
		},

		{
			input: `{
"name": "frontmatter",
}

rest of the file`,
			expectedFunc: expMatterErr,
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
			expectedFunc: expValidMatter,
			formats: []*frontmatter.Format{
				frontmatter.NewFormat("...", "...", yaml.Unmarshal),
			},
		},
	}

	failFunc := func(exp, act interface{}) {
		t.Fatalf("Not equal: \n"+
			"expected: %v\n"+
			"actual  : %v", exp, act)
	}

	checkFunc := func(mExp, mAct *matter, rExp, rAct string, eExp, eAct bool) {
		if eExp != eAct {
			failFunc(eExp, eAct)
		}
		if mExp.Name != mAct.Name {
			failFunc(mExp.Name, mAct.Name)
		}
		if strings.Join(mExp.Tags, ",") != strings.Join(mAct.Tags, ",") {
			failFunc(mExp.Tags, mAct.Tags)
		}
		if mExp.Metadata.Size != mAct.Metadata.Size {
			failFunc(mExp.Metadata.Size, mAct.Metadata.Size)
		}
		if rExp != rAct {
			failFunc(rExp, rAct)
		}
	}

	for _, tc := range testCases {
		input := tc.input

		// Get expected data.
		mExp, rExp, hasErr := tc.expectedFunc(input)

		// Parse reader.
		mAct := &matter{}
		rest, err := frontmatter.Parse(strings.NewReader(input), mAct, tc.formats...)
		checkFunc(mExp, mAct, rExp, string(rest), hasErr, err != nil)
	}
}
