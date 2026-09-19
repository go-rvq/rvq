package admin

import (
	"encoding/json"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// exampleConfig is a singleton whose whole configuration lives in one JSON map
// (Data). It stands in, in this documentation example, for a real config model
// such as an application's SEO/site settings — history versions the payload and
// the comparator/partial-revert show and edit it as YAML.
type exampleConfig struct {
	ID   uint
	Data map[string]any
}

// ExampleYAMLConfigHistory shows how to activate git-like history on a model
// whose versioned field is a JSON map, and adapt the comparator and partial
// revert to work on it as YAML — the same wiring an application uses for a
// settings/config model. It is a documentation example (see the history README);
// db and mb are the app's gorm DB and the model builder.
//
// @snippet_begin(HistoryYAMLConfigExample)
func ExampleYAMLConfigHistory(db *gorm.DB, mb *presets.ModelBuilder) {
	// Version the Data field. History records a revision on every save, browsable
	// and revertible under /<model>/revisions.
	mh := New(db).
		Model(mb).
		Fields("Data").
		Build()

	// Show the JSON payload as YAML in the comparator: line-by-line, with the
	// changed content highlighted within each line (GoLand/IntelliJ style) and
	// Prism syntax coloring.
	mh.FieldDiffHandler("Data", func(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
		oldC, newC = PrismDiffComponents("yaml",
			exampleJSONToYAML(in.OldSnap["Data"]), exampleJSONToYAML(in.NewSnap["Data"]))
		return oldC, newC, nil, true
	})

	// A JSON map is not a string, so it joins the text-based partial revert
	// through a content codec: it is diffed/edited as YAML text (ToText) and, on
	// apply, the patched YAML is parsed back into the map (Apply). Registering a
	// codec is what makes a structured field accept partial revert.
	mh.FieldContentHandler("Data", "yaml", FieldContentCodec{
		ToText: func(rawJSON string) string { return exampleJSONToYAML([]byte(rawJSON)) },
		Apply: func(obj any, text string) error {
			c := obj.(*exampleConfig)
			data := map[string]any{}
			if err := yaml.Unmarshal([]byte(text), &data); err != nil {
				return err
			}
			c.Data = data
			return nil
		},
	})
}

// @snippet_end

// exampleJSONToYAML converts a field's JSON payload to YAML for display and
// patching (an empty or null payload becomes "{}").
func exampleJSONToYAML(raw []byte) string {
	if len(raw) == 0 || string(raw) == "null" {
		return "{}"
	}
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return string(raw)
	}
	b, err := yaml.Marshal(v)
	if err != nil {
		return string(raw)
	}
	return string(b)
}
