package players

import (
	"github.com/COMTOP1/api/services/afc/players"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repo struct {
	players *players.Store
	access  *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		players: players.NewStore(scope),
		access:  access,
	}
}
