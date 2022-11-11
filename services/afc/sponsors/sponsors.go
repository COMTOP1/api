package sponsors

import "github.com/couchbase/gocb/v2"

type (
	SponsorsRepo interface {
	}

	Store struct {
		scope *gocb.Scope
	}

	Sponsor struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
