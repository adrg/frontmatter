package frontmatter

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Format describes a front matter. It holds all the information
// necessary in order to detect and decode a front matter format.
type Format struct {
	// Start defines the starting delimiter of the front matter.
	// E.g.: `---` or `---yaml`.
	Start string

	// End defines the ending delimiter of the front matter.
	// E.g.: `---`.
	End string

	// Unmarshal defines the unmarshal function used to decode
	// the front matter data, after it has been detected.
	// E.g.: json.Unmarshal (from the `encoding/json` package).
	Unmarshal func(data []byte, v interface{}) error

	// UnmarshalDelims specifies whether the front matter
	// delimiters are included in the data to be unmarshaled.
	// Should be `false` in most cases.
	UnmarshalDelims bool

	// RequiresNewLine specifies whether a new (empty) line is
	// required after the front matter.
	// Should be `false` in most cases.
	RequiresNewLine bool
}

func newFormat(start, end string, unmarshal func([]byte, interface{}) error,
	unmarshalDelims, requiresNewLine bool) *Format {
	return &Format{
		Start:           start,
		End:             end,
		Unmarshal:       unmarshal,
		UnmarshalDelims: unmarshalDelims,
		RequiresNewLine: requiresNewLine,
	}
}

func defaultFormats() []*Format {
	return []*Format{
		// YAML.
		newFormat("---", "---", yaml.Unmarshal, false, false),
		newFormat("---yaml", "---", yaml.Unmarshal, false, false),
		// TOML.
		newFormat("+++", "+++", toml.Unmarshal, false, false),
		newFormat("---toml", "---", toml.Unmarshal, false, false),
		// JSON.
		newFormat(";;;", ";;;", json.Unmarshal, false, false),
		newFormat("---json", "---", json.Unmarshal, false, false),
		newFormat("{", "}", json.Unmarshal, true, true),
	}
}
