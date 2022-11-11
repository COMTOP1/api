package teams

import (
	"github.com/COMTOP1/api/services/afc/teams"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repo struct {
	teams  *teams.Store
	access *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		teams:  teams.NewStore(scope),
		access: access,
	}
}
