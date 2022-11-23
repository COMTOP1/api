package programmes

import (
	"github.com/COMTOP1/api/services/afc/programmes"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
    "github.com/labstack/echo/v4"
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

func (r *Repo) ListAllProgrammes(c echo.Context) error {
    return nil
}