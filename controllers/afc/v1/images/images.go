package images

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/images"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

func (r *Repo) GetImageById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("GetImageById failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.images.GetImageById(id)
	if err != nil {
		err = fmt.Errorf("GetImageById failed to get image: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllImages(c echo.Context) error {
	p, err := r.images.ListAllImages()
	if err != nil {
		err = fmt.Errorf("ListAllImages failed to get all image: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) AddImage(c echo.Context) error {
	var p *images.Image
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("AddImage failed to bind image: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.images.AddImage(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) DeleteImage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	_, err = r.images.GetImageById(id)
	if err != nil {
		err = fmt.Errorf("DeleteImage failed to get image: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.images.DeleteImage(id)
	if err != nil {
		err = fmt.Errorf("DeleteImage failed to delete image: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
