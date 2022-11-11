package affiliations

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	AffiliationsRepo interface {
		GetAffiliation(ctx context.Context, id uint) (Affiliation, error)
	}

	Store struct {
		scope *gocb.Scope
	}

	Affiliation struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
