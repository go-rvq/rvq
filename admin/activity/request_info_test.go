package activity

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestRequestInfoFromRequest(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(r *httptest.ResponseRecorder) map[string]string
		headers map[string]string
		remote  string
		wantIP  string
	}{
		{name: "remote addr with port", remote: "10.1.2.3:5000", wantIP: "10.1.2.3"},
		{name: "x-forwarded-for single", remote: "10.0.0.1:1", headers: map[string]string{"X-Forwarded-For": "200.1.2.3"}, wantIP: "200.1.2.3"},
		{name: "x-forwarded-for chain", remote: "10.0.0.1:1", headers: map[string]string{"X-Forwarded-For": "200.1.2.3, 10.0.0.2"}, wantIP: "200.1.2.3"},
		{name: "x-real-ip", remote: "10.0.0.1:1", headers: map[string]string{"X-Real-Ip": "201.4.5.6"}, wantIP: "201.4.5.6"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = c.remote
			r.Header.Set("User-Agent", "TestBrowser/1.0")
			for k, v := range c.headers {
				r.Header.Set(k, v)
			}
			info := RequestInfoFromRequest(r)
			if info.IP != c.wantIP {
				t.Errorf("IP = %q, want %q", info.IP, c.wantIP)
			}
			if info.UserAgent != "TestBrowser/1.0" {
				t.Errorf("UserAgent = %q", info.UserAgent)
			}
		})
	}
}

func TestContextWithRequestInfo(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.1.2.3:5000"
	r.Header.Set("User-Agent", "TestBrowser/1.0")

	ctx := ContextWithRequestInfo(context.Background(), r)
	info := RequestInfoFromContext(ctx)
	if info == nil {
		t.Fatal("info not stored in context")
	}
	if info.IP != "10.1.2.3" || info.UserAgent != "TestBrowser/1.0" {
		t.Errorf("info = %+v", info)
	}

	if RequestInfoFromContext(context.Background()) != nil {
		t.Error("empty context must return nil info")
	}
}

func TestActivityLogRequestInfoSetter(t *testing.T) {
	var log ActivityLogInterface = &ActivityLog{}
	s, ok := log.(RequestInfoSetter)
	if !ok {
		t.Fatal("ActivityLog must implement RequestInfoSetter")
	}
	s.SetIP("10.0.0.9")
	s.SetUserAgent("UA")
	al := log.(*ActivityLog)
	if al.IP != "10.0.0.9" || al.UserAgent != "UA" {
		t.Errorf("log = %+v", al)
	}
}
