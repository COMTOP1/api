package affiliations

import (
	"github.com/COMTOP1/api/services/afc/affiliations"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
    "github.com/labstack/echo/v4"
)

// Repo stores our dependencies
type Repo struct {
	affiliations *affiliations.Store
	access       *utils.Accesser
}

// NewRepo creates our data store
func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		affiliations: affiliations.NewStore(scope),
		access:       access,
	}
}

func (r *Repo) ListAllAffiliations(c echo.Context) error {
    return nil
}