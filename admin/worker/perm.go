package worker

import "net/http"

// examples:.
// permPolicy.On("*")
// permPolicy.On("workers:upload_posts")
const (
	PermEdit = "@edit"
)

func editIsAllowed(r *http.Request, jobName string) error {
	return permVerifier.Do(PermEdit).SnakeOn(jobName).WithReq(r).IsAllowed()
}
