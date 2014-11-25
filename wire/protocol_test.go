package wire

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPeriodMarshalEmpty(t *testing.T) {
	type foo struct {
		When Period `json:",omitempty"`
	}

	b, err := json.Marshal(foo{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expect := []byte("{}")
	if !bytes.Equal(expect, b) {
		t.Fatalf("expected: %q, got: %q", expect, b)
	}
}
