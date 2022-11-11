package programmes

import (
	"github.com/COMTOP1/api/services/afc/programmes"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repo struct {
	programmes *programmes.Store
	access     *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		programmes: programmes.NewStore(scope),
		access:     access,
	}
}
