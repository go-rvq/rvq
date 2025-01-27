package editorjs

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrEmbServiceNotSupported = errors.New("Only Youtube and Vime Embeds are supported right now.")
)

type BlockTypeNotSupported struct {
	BlockType string
}

func (e *BlockTypeNotSupported) Error() string {
	return fmt.Sprintf("Block type %q is not supported.", e.BlockType)
}
