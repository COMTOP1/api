package whatsOn

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/whatsOn"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Repo struct {
	whatsOn    *whatsOn.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		whatsOn:    whatsOn.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) GetWhatsOnByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("GetWhatsOnByID failed to get id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	w, err := r.whatsOn.GetWhatsOn(id)
	if err != nil {
		err = fmt.Errorf("GetWhatsOnByID failed to get whatson: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) GetWhatsOnLatest(c echo.Context) error {
	w, err := r.whatsOn.GetWhatsOnLatest()
	if err != nil {
		err = fmt.Errorf("GetWhatsOnByID failed to get whatson: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) ListAllWhatsOn(c echo.Context) error {
	w, err := r.whatsOn.ListAllWhatsOn()
	if err != nil {
		err = fmt.Errorf("ListAllWhatsOn failed to get all whatsOn: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) AddWhatsOn(c echo.Context) error {
	var w *whatsOn.WhatsOn
	err := c.Bind(w)
	if err != nil {
		err = fmt.Errorf("AddWhatsOn failed to get whatsOn: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.whatsOn.AddWhatsOn(w)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) EditWhatsOn(c echo.Context) error {
	var w *whatsOn.WhatsOn
	err := c.Bind(w)
	if err != nil {
		err = fmt.Errorf("EditWhatsOn failed to get whatsOn: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.whatsOn.EditWhatsOn(w)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) DeleteWhatsOn(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	_, err = r.whatsOn.GetWhatsOn(id)
	if err != nil {
		err = fmt.Errorf("DeleteUser failed to get user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.whatsOn.DeleteWhatsOn(id)
	if err != nil {
		err = fmt.Errorf("DeleteUser failed to delete user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
