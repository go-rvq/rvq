package presets

import (
	"testing"

	"github.com/go-rvq/rvq/admin/model"
)

type parentsReply struct {
	ID   uint
	Body string
}

// Under two parents: the ids of the parents the outermost first (as the
// request carries them), each on its own model.
func TestPermParentsOrder(t *testing.T) {
	_, posts, comments := uniqueApp(nil)
	replies := posts.Builder().Model(&parentsReply{}, ModelWithID("replies"), ModelNotInMenu())
	comments.AddChild(replies)
	id := func(_ *ModelBuilder, v string) model.ID { return model.NewID(v) }
	v := replies.Permissioner().Verifier(id(replies, "3"), id(posts, "1"), id(comments, "2"))
	if r := v.PreferredResource(); r != "admin:posts:<1>:comments:<2>:replies:<3>:" {
		t.Errorf("the record: %q", r)
	}
	if r := replies.Permissioner().ListVerifier(id(posts, "1"), id(comments, "2")).PreferredResource(); r != "admin:posts:<1>:comments:<2>:replies:" {
		t.Errorf("the listing: %q", r)
	}
}
