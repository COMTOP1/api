package news

import (
	"github.com/COMTOP1/api/services/afc/news"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repo struct {
	news   *news.Store
	access *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		news:   news.NewStore(scope),
		access: access,
	}
}
