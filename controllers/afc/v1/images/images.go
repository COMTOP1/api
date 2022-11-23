package images

import (
	"github.com/COMTOP1/api/services/afc/images"
	"github.com/COMTOP1/api/utils"
    "github.com/couchbase/gocb/v2"
    "github.com/labstack/echo/v4"
)

type Repo struct {
	images *images.Store
	access *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		images: images.NewStore(scope),
		access: access,
	}
}

func (r *Repo) ListAllImages(c echo.Context) error {
    return nil
}