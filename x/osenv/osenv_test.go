package osenv

import "testing"

// A variable is its value, else its default; it is recorded, where it was read.
func TestGet(t *testing.T) {
	t.Setenv("OSENV_TEST_A", "x")
	t.Setenv("OSENV_TEST_B", "")
	t.Setenv("OSENV_TEST_C", "7")
	if v := Get("OSENV_TEST_A", "a", "d"); v != "x" {
		t.Errorf("A = %q", v)
	}
	if v := Get("OSENV_TEST_B", "b", "d"); v != "d" {
		t.Errorf("B = %q", v)
	}
	if v := GetInt64("OSENV_TEST_C", "c", 1); v != 7 {
		t.Errorf("C = %d", v)
	}
	if v := GetBool("OSENV_TEST_D", "d", true); !v {
		t.Errorf("D = %v", v)
	}
	vs := Vars()
	if len(vs) < 4 || vs[len(vs)-4].Key != "OSENV_TEST_A" || vs[len(vs)-4].At != "github.com/go-rvq/rvq/x/osenv.TestGet" {
		t.Errorf("Vars: %+v", vs)
	}
	if s := vs[len(vs)-3].String(); s != "OSENV_TEST_B: b, (default: d), at: github.com/go-rvq/rvq/x/osenv.TestGet" {
		t.Errorf("String: %q", s)
	}
}
