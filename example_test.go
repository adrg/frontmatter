package frontmatter_test

import (
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

func ExampleParse_yAML() {
	r := strings.NewReader(`
---
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
---
rest of the content`)

	var matter struct {
		Name string   `yaml:"name" toml:"name" json:"name"`
		Tags []string `yaml:"tags" toml:"tags" json:"tags"`
	}

	rest, err := frontmatter.Parse(r, &matter)
	if err != nil {
		// Treat error.
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	fmt.Printf("%+v\n", matter)
	fmt.Println(string(rest))

	// Output:
	// {Name:frontmatter Tags:[go yaml json toml]}
	// rest of the content
}

func ExampleParse_tOML() {
	r := strings.NewReader(`
+++
name = "frontmatter"
tags = ["go", "yaml", "json", "toml"]
+++
rest of the content`)

	var matter struct {
		Name string   `yaml:"name" toml:"name" json:"name"`
		Tags []string `yaml:"tags" toml:"tags" json:"tags"`
	}

	rest, err := frontmatter.Parse(r, &matter)
	if err != nil {
		// Treat error.
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	fmt.Printf("%+v\n", matter)
	fmt.Println(string(rest))

	// Output:
	// {Name:frontmatter Tags:[go yaml json toml]}
	// rest of the content
}

func ExampleParse_jSON() {
	r := strings.NewReader(`
{
  "name": "frontmatter",
  "tags": ["go", "yaml", "json", "toml"]
}

rest of the content`)

	var matter struct {
		Name string   `yaml:"name" toml:"name" json:"name"`
		Tags []string `yaml:"tags" toml:"tags" json:"tags"`
	}

	rest, err := frontmatter.Parse(r, &matter)
	if err != nil {
		// Treat error.
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	fmt.Printf("%+v\n", matter)
	fmt.Println(string(rest))

	// Output:
	// {Name:frontmatter Tags:[go yaml json toml]}
	// rest of the content
}

func ExampleParse_custom() {
	r := strings.NewReader(`
...
name: "frontmatter"
tags: ["go", "yaml", "json", "toml"]
...
rest of the content`)

	var matter struct {
		Name string   `yaml:"name"`
		Tags []string `yaml:"tags"`
	}

	formats := []*frontmatter.Format{
		frontmatter.NewFormat("...", "...", yaml.Unmarshal),
	}

	rest, err := frontmatter.Parse(r, &matter, formats...)
	if err != nil {
		// Treat error.
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	fmt.Printf("%+v\n", matter)
	fmt.Println(string(rest))

	// Output:
	// {Name:frontmatter Tags:[go yaml json toml]}
	// rest of the content
}
