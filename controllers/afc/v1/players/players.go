package players

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/players"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	players    *players.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		players:    players.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListAllPlayersByTeamID(c echo.Context) error {
	return nil
}
