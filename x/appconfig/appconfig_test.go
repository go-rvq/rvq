package appconfig

import "testing"

func TestSection(t *testing.T) {
	Reset()
	defer Reset()
	Load([]byte("env: {A: b}\nfstools:\n  terminal:\n    sandbox: landlock\n"))
	var c struct {
		Terminal struct {
			Sandbox string `yaml:"sandbox"`
			Idle    string `yaml:"idle"`
		} `yaml:"terminal"`
	}
	c.Terminal.Idle = "30m"
	if ok, err := Section("fstools", &c); !ok || err != nil || c.Terminal.Sandbox != "landlock" || c.Terminal.Idle != "30m" {
		t.Errorf("%v %v %+v", ok, err, c)
	}
	if ok, _ := Section("none", &c); ok {
		t.Error("none found")
	}
}
