package players

import (
	"github.com/COMTOP1/api/services/afc/players"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
    "github.com/labstack/echo/v4"
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

func (r *Repo) ListAllPlayersByTeamID(c echo.Context) error {
    return nil
}