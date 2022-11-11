package images

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	ImagesRepo interface {
		GetImage(ctx context.Context, id uint) (Image, error)
	}

	Store struct {
		scope *gocb.Scope
	}

	Image struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
