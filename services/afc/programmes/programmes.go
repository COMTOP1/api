package programmes

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	ProgrammesRepo interface {
		GetProgramme(ctx context.Context, id uint) (Programme, error)
	}

	Store struct {
		scope *gocb.Scope
	}

	Programme struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
