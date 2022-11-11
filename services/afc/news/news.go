package news

import "github.com/couchbase/gocb/v2"

type (
	NewsRepo interface {
	}

	Store struct {
		scope *gocb.Scope
	}

	News struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
