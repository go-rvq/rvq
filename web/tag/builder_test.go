package tag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagBuilder_RemoveAttr(t *testing.T) {
	type myTag struct {
		TagBuilder[*myTag]
	}
	tag := NewTag(&myTag{}, "my-tag").Attr("a", "b").Attr("c", "d")

	toS := func() string {
		b, _ := tag.MarshalHTML(context.Background())
		return string(b)
	}

	require.Equal(t, "\n<my-tag a='b' c='d'></my-tag>\n", toS())
	tag.RemoveAttr("a")
	require.Equal(t, "\n<my-tag c='d'></my-tag>\n", toS())
	tag.RemoveAttr("c")
	require.Equal(t, "\n<my-tag></my-tag>\n", toS())
}
