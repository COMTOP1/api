package teams

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	TeamsRepo interface {
		GetTeam(ctx context.Context, id uint) (Team, error)
	}

	Store struct {
		scope *gocb.Scope
	}

	Team struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
