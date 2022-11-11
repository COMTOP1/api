package admin

import (
	"context"
	"github.com/couchbase/gocb/v2"
)

type (
	AdminRepo interface {
		GetJWT(ctx context.Context, email string) (JWTToken, error)
		GetSSOJWT(ctx context.Context) (JWTToken, error)
	}

	// Store contains our dependency
	Store struct {
		scope *gocb.Scope
	}
)

type (
	JWTToken struct {
		JWTToken string `json:"jwt_token"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}
