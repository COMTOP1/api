package documents

import (
	"github.com/COMTOP1/api/services/afc/documents"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repo struct {
	documents *documents.Store
	access    *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		documents: documents.NewStore(scope),
		access:    access,
	}
}
