package images

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/images"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	images     *images.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		images:     images.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListAllImages(c echo.Context) error {
	return nil
}
