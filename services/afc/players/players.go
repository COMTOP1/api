package players

import "github.com/couchbase/gocb/v2"

type (
	PlayersRepo interface {
	}

	Store struct {
		scope *gocb.Scope
	}

	Player struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
