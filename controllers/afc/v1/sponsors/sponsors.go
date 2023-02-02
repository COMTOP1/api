package sponsors

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/sponsors"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	sponsors   *sponsors.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		sponsors:   sponsors.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListALlSponsors(c echo.Context) error {
	return nil
}

func (r *Repo) ListAllSponsorsByTeamID(c echo.Context) error {
	return nil
}
