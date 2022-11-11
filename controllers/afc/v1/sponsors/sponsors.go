package sponsors

import (
	"github.com/COMTOP1/api/services/afc/sponsors"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repo struct {
	sponsors *sponsors.Store
	access   *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		sponsors: sponsors.NewStore(scope),
		access:   access,
	}
}
