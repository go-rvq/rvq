package schemaform

import (
	"strings"
	"testing"
)

// The fields are checked by the contracts of their types (contracts/
// contracts.gad, linked into each schema as the module schemaform): an enum's
// options are its members, by name or value; a select of gad's values has its
// options.
func TestContractsCheckFields(t *testing.T) {
	b := New()
	if _, err := b.Parse(`enum Service { Kitchen, Bathroom }
export class default {
    service Service
    [options=["Kitchen", 2]] some Service
    [options=["a", "b"]] pick str
}`); err != nil {
		t.Fatal(err)
	}
	_, err := b.Parse(`enum Service { Kitchen, Bathroom }
export class default {
    [options=["x", 9]] bad Service
}`)
	if err == nil || !strings.Contains(err.Error(), "bad: meta \"options\": item 0: ErrIncompatibleAssign: ‹x› is not a member of ‹Service›") ||
		!strings.Contains(err.Error(), "item 1: ErrIncompatibleAssign: ‹9› is not a member of ‹Service›") {
		t.Fatalf("an enum's options not its members: %v", err)
	}
	// the plugin of the contracts, in the schema itself
	if _, err := b.Parse(`enum Service { Kitchen, Bathroom }
class Probe { service Service }
if schemaform.field(Probe.service).contract.@name != "TEnumField" { throw "not the enum's contract" }
export class default { name str }`); err != nil {
		t.Fatal(err)
	}
}
