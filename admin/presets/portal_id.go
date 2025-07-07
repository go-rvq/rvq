package presets

import (
	"fmt"
	"net/http"
	"time"
)

func GetPortalID(r *http.Request) (v string) {
	return r.FormValue(ParamPortalID)
}

func GetOrNewPortalID(r *http.Request) (v string) {
	if v = r.FormValue(ParamPortalID); v == "" {
		v = fmt.Sprintf("_%d", time.Now().UnixNano())
		r.Form.Set(ParamPortalID, v)
	}
	return
}
