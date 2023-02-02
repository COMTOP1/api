package programmes

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/programmes"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	programmes *programmes.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		programmes: programmes.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListAllProgrammes(c echo.Context) error {
	return nil
}
