package documents

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/documents"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	documents  *documents.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		documents:  documents.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListAllDocuments(c echo.Context) error {
	return nil
}
