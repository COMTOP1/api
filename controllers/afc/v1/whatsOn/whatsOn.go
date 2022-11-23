package whatsOn

import (
	"github.com/COMTOP1/api/services/afc/whatsOn"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
    "github.com/labstack/echo/v4"
)

type Repo struct {
	whatsOn *whatsOn.Store
	access  *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		whatsOn: whatsOn.NewStore(scope),
		access:  access,
	}
}

func (r *Repo) ListAllWhatsOn(c echo.Context) error {
    return nil
}

func (r *Repo) GetWhatsOnByID(c echo.Context) error {
    return nil
}