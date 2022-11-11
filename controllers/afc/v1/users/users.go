package users

import (
	"github.com/COMTOP1/api/services/afc/users"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

// Repo stores our dependencies
type Repo struct {
	users  *users.Store
	access *utils.Accesser
}

// NewRepo creates our data store
func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		users:  users.NewStore(scope),
		access: access,
	}
}
