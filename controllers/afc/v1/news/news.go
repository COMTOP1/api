package news

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/news"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
)

type Repo struct {
	news       *news.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		news:       news.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) ListAllNews(c echo.Context) error {
	return nil
}

func (r *Repo) GetNewsByID(c echo.Context) error {
	return nil
}

func (r *Repo) GetNewsLatest(c echo.Context) error {
	return nil
}
