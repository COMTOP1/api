package affiliations

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/affiliations"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

// Repo stores our dependencies
type Repo struct {
	affiliations *affiliations.Store
	controller   controllers.Controller
}

// NewRepo creates our data store
func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		affiliations: affiliations.NewStore(scope),
		controller:   controller,
	}
}

func (r *Repo) ListAllAffiliations(c echo.Context) error {
	return nil
}
