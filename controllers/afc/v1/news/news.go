package news

import (
	"github.com/COMTOP1/api/services/afc/news"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
    "github.com/labstack/echo/v4"
)

type Repo struct {
	news   *news.Store
	access *utils.Accesser
}

func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		news:   news.NewStore(scope),
		access: access,
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
