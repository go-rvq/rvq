package tag

import (
	"context"
	h "github.com/go-rvq/htmlgo"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagBuilder_RemoveAttr(t *testing.T) {
	type myTag struct {
		TagBuilder[*myTag]
	}
	tag := NewTag(&myTag{}, "my-tag").Attr("a", "b").Attr("c", "d")

	// htmlgo.Marshal renders the component as-is; the builder method that used
	// to wrap it in newlines is gone
	toS := func() string {
		b, _ := h.Marshal(tag, context.Background())
		return string(b)
	}

	require.Equal(t, "<my-tag a='b' c='d'></my-tag>", toS())
	tag.RemoveAttr("a")
	require.Equal(t, "<my-tag c='d'></my-tag>", toS())
	tag.RemoveAttr("c")
	require.Equal(t, "<my-tag></my-tag>", toS())
}
