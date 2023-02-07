package programmes

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/programmes"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"unicode"
)

type Repo struct {
	programmes *programmes.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		programmes: programmes.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) GetProgrammeById(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("GetProgrammeById failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.programmes.GetProgrammeById(id)
	if err != nil {
		err = fmt.Errorf("GetProgrammeById failed to get programme: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllProgrammes(c echo.Context) error {
	p, err := r.programmes.ListAllProgrammes()
	if err != nil {
		err = fmt.Errorf("ListAllProgrammes failed to get all programme: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) AddProgramme(c echo.Context) error {
	var p *programmes.Programme
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("AddProgramme failed to bind programme: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.programmes.AddProgramme(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) DeleteProgramme(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	_, err = r.programmes.GetProgrammeById(id)
	if err != nil {
		err = fmt.Errorf("DeleteProgramme failed to get programme: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.programmes.DeleteProgramme(id)
	if err != nil {
		err = fmt.Errorf("DeleteProgramme failed to delete programme: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
