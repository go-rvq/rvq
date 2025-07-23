package presets

import (
	"iter"

	"github.com/go-rvq/rvq/x/perm"
)

func (mb *ModelBuilder) AllVerifiers() iter.Seq[*perm.PermVerifierBuilder] {
	return func(yield func(*perm.PermVerifierBuilder) bool) {
		for _, v := range mb.verifiers {
			if !yield(v) {
				return
			}
		}
	}
}

func (d *DetailingBuilder) AllVerifiers() iter.Seq[*perm.PermVerifierBuilder] {
	return func(yield func(*perm.PermVerifierBuilder) bool) {
		for _, v := range d.verifiers {
			if !yield(v) {
				return
			}
		}

		for v := range d.pagesRegistrator.BuildedVerifiers() {
			if !yield(v) {
				return
			}
		}
	}
}
