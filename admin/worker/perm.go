package worker

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/presets"
)

// The jobs are a model of the admin (scope presets): doing a job of a kind —
// creating it, rerunning it, aborting it, updating it — is an action of the
// model, "!" and the job's name:
//
//	:system/:jobs:!upload_posts
//
// — through its groups, or by its unique name (:jobs:…); "*" for
// any kind.

// jobsModel is the model of the jobs (Install).
var jobsModel *presets.ModelBuilder

func editIsAllowed(r *http.Request, jobName string) error {
	return jobsModel.Permissioner().ReqListActioner(r, presets.ActionPerm(jobName)).IsAllowed()
}
