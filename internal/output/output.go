package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Format is the output presentation format
type Format string

const (
	FormatText  Format = "text"
	FormatHuman Format = "human"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Renderer coordinates rendering results structures
type Renderer interface {
	Render(w io.Writer, data interface{}) error
}

// ConfiguredRenderer handles standard CLI output formatting
type ConfiguredRenderer struct {
	Format      Format
	SilentMode  bool
	VerboseMode bool
}

func NewConfiguredRenderer(fmt Format, silent, verbose bool) *ConfiguredRenderer {
	return &ConfiguredRenderer{
		Format:      fmt,
		SilentMode:  silent,
		VerboseMode: verbose,
	}
}

func (r *ConfiguredRenderer) Render(w io.Writer, data interface{}) error {
	if r.SilentMode {
		return nil
	}

	switch r.Format {
	case FormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case FormatYAML:
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		return encoder.Encode(data)
	default:
		// Human readable text presentation
		if str, ok := data.(string); ok {
			fmt.Fprintln(w, str)
			return nil
		}
		if stringer, ok := data.(fmt.Stringer); ok {
			fmt.Fprintln(w, stringer.String())
			return nil
		}
		// Fallback to formatted printed text
		fmt.Fprintf(w, "%+v\n", data)
		return nil
	}
}

// Standard Console Dumps helpers
var Stdout io.Writer = os.Stdout
var Stderr io.Writer = os.Stderr
