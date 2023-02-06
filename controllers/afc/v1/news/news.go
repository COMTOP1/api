package news

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/news"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

func (r *Repo) GetNewsById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("GetNewsById failed to get id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	n, err := r.news.GetNewsById(id)
	if err != nil {
		err = fmt.Errorf("GetNewsById failed to get news: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (r *Repo) GetNewsLatest(c echo.Context) error {
	n, err := r.news.GetNewsLatest()
	if err != nil {
		err = fmt.Errorf("GetNewsById failed to get news: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (r *Repo) ListAllNews(c echo.Context) error {
	n, err := r.news.ListAllNews()
	if err != nil {
		err = fmt.Errorf("ListAllNews failed to get all news: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (r *Repo) AddNews(c echo.Context) error {
	var n *news.News
	err := c.Bind(&n)
	if err != nil {
		err = fmt.Errorf("AddNews failed to bind news: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.news.AddNews(n)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (r *Repo) EditNews(c echo.Context) error {
	var n *news.News
	err := c.Bind(&n)
	if err != nil {
		err = fmt.Errorf("EditNews failed to get news: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.news.EditNews(n)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (r *Repo) DeleteNews(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	_, err = r.news.GetNewsById(id)
	if err != nil {
		err = fmt.Errorf("DeleteNews failed to get news: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.news.DeleteNews(id)
	if err != nil {
		err = fmt.Errorf("DeleteNews failed to delete news: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
