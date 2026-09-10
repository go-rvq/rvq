package models

import (
	"bytes"
	"testing"
)

func TestHashStringHex(t *testing.T) {
	h := Hash{0xab, 0x12, 0xff, 0x00}
	if got, want := h.String(), "ab12ff00"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestHashParseRoundTrip(t *testing.T) {
	h := Hash{0xde, 0xad, 0xbe, 0xef}
	got, err := Hash{}.Parse(h.String())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if !bytes.Equal(got, h) {
		t.Fatalf("Parse(String()) = %x, want %x", got, h)
	}
}

func TestHashParseBad(t *testing.T) {
	if _, err := (Hash{}).Parse("zz"); err == nil {
		t.Fatal("Parse of non-hex should error")
	}
}

func TestHashValueScanRoundTrip(t *testing.T) {
	h := Hash{1, 2, 3, 250}

	v, err := h.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	b, ok := v.([]byte)
	if !ok {
		t.Fatalf("Value type = %T, want []byte", v)
	}

	var back Hash
	if err := back.Scan(b); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !bytes.Equal(back, h) {
		t.Fatalf("Scan(Value()) = %x, want %x", back, h)
	}
}

func TestHashScanString(t *testing.T) {
	var h Hash
	if err := h.Scan("\x01\x02"); err != nil {
		t.Fatalf("Scan string: %v", err)
	}
	if !bytes.Equal(h, Hash{1, 2}) {
		t.Fatalf("Scan(string) = %x, want 0102", h)
	}
}

func TestHashValueEmptyIsNil(t *testing.T) {
	v, err := Hash(nil).Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if v != nil {
		t.Fatalf("empty Value() = %v, want nil", v)
	}
}
