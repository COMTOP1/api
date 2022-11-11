package whatsOn

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	WhatsOnRepo interface {
		GetWhatsOn(ctx context.Context, id uint) (WhatsOn, error)
	}

	Store struct {
		scope *gocb.Scope
	}

	WhatsOn struct {
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
