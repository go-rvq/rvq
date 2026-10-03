package worker

import "net/http"

// The scope workers: the jobs, each by its name — editing one, a job that
// takes arguments, is "workers:<name>:@edit".
//
// examples:
// permPolicy.On("*")
// permPolicy.On("workers:upload_posts:*")
const (
	PermEdit = "@edit"
)

func editIsAllowed(r *http.Request, jobName string) error {
	return permVerifier.Do(PermEdit).SnakeOn(jobName).WithReq(r).IsAllowed()
}
