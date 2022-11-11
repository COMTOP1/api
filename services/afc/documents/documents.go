package documents

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	DocumentsRepo interface {
		GetDocument(ctx context.Context, id uint) (Document, error)
	}

	Store struct {
		scope *gocb.Scope
	}

	Document struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
