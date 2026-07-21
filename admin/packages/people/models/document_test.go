package models

import "testing"

func TestValidCPF(t *testing.T) {
	valid := []string{"529.982.247-25", "52998224725"}
	for _, d := range valid {
		if !ValidCPF(d) {
			t.Errorf("ValidCPF(%q) = false, want true", d)
		}
	}
	invalid := []string{"111.111.111-11", "529.982.247-24", "123", ""}
	for _, d := range invalid {
		if ValidCPF(d) {
			t.Errorf("ValidCPF(%q) = true, want false", d)
		}
	}
}

func TestValidCNPJ(t *testing.T) {
	if !ValidCNPJ("11.222.333/0001-81") {
		t.Error("ValidCNPJ(11.222.333/0001-81) = false, want true")
	}
	for _, d := range []string{"11.222.333/0001-80", "00000000000000", "abc"} {
		if ValidCNPJ(d) {
			t.Errorf("ValidCNPJ(%q) = true, want false", d)
		}
	}
}
