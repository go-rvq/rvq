package web

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
)

var enc = base64.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_").WithPadding(base64.NoPadding)

type Callback struct {
	Scripts       []string
	ReloadPortals []string
}

func (c *Callback) UnmarshalJSON(i []byte) (err error) {
	c.Scripts = nil
	c.ReloadPortals = nil
	var encoded string
	if err = json.Unmarshal(i, &encoded); err != nil {
		return
	}
	c.Decode(encoded)
	return
}

func (c *Callback) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Callback) AddScript(script string) *Callback {
	c.Scripts = append(c.Scripts, script)
	return c
}

func (c *Callback) AddReloadPortal(portal string) *Callback {
	c.ReloadPortals = append(c.ReloadPortals, portal)
	return c
}

func (c *Callback) Encode() string {
	if len(c.ReloadPortals) == 0 && len(c.Scripts) == 0 {
		return ""
	}

	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.BestCompression)
	w.Write([]byte(strings.Join(c.ReloadPortals, ",")))
	w.Write([]byte{';'})
	for _, script := range c.Scripts {
		w.Write([]byte(script))
		w.Write([]byte{';'})
	}
	w.Close()
	return enc.EncodeToString(buf.Bytes())
}

func (c *Callback) Decode(s string) {
	if len(s) == 0 {
		*c = Callback{}
		return
	}
	if data, err := enc.DecodeString(s); err == nil {
		r := flate.NewReader(bytes.NewBuffer(data))
		var b bytes.Buffer
		io.Copy(&b, r)
		r.Close()
		data = b.Bytes()
		pos := bytes.IndexByte(data, ';')
		if pos > 0 {
			c.ReloadPortals = strings.Split(string(data[:pos]), ",")
		}
		data = data[pos+1:]
		if len(data) > 1 {
			c.Scripts = []string{string(data)}
		}
	}
}

func DecodeCallback(s string) (c Callback) {
	if s != "" {
		c.Decode(s)
	}
	return
}

func (c *Callback) Script() string {
	return strings.Join(c.Scripts, ";")
}

func (c *Callback) String() string {
	return c.Encode()
}

func CallbackScript(script string) (c *Callback) {
	c = &Callback{}
	c.AddScript(script)
	return
}
