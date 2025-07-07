package web

import "net/http"

const ExecuteEventParam = "__execute_event__"

func CurrentExecuteEvent(r *http.Request) string {
	return r.Form.Get(ExecuteEventParam)
}
