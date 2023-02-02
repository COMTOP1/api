package teams

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/teams"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	teams      *teams.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		teams:      teams.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListAllTeams(c echo.Context) error {
	return nil
}

func (r *Repo) GetTeamByID(c echo.Context) error {
	return nil
}

func (r *Repo) GetTeamManagerByID(c echo.Context) error {
	return nil
}
