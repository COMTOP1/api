package programmeSeasons

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/programmeSeasons"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Repo struct {
	programmeSeasons *programmeSeasons.Store
	controller       controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		programmeSeasons: programmeSeasons.NewStore(scope),
		controller:       controller,
	}
}

func (r *Repo) GetProgrammeSeasonById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("GetProgrammeSeasonById failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.programmeSeasons.GetProgrammeSeasonById(id)
	if err != nil {
		err = fmt.Errorf("GetProgrammeSeasonById failed to get programme season: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllProgrammeSeasons(c echo.Context) error {
	p, err := r.programmeSeasons.ListAllProgrammeSeasons()
	if err != nil {
		err = fmt.Errorf("ListAllProgrammeSeasons failed to get all programme season: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) AddProgrammeSeason(c echo.Context) error {
	var p *programmeSeasons.ProgrammeSeason
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("AddProgrammeSeason failed to bind programme season: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.programmeSeasons.AddProgrammeSeason(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) EditProgrammeSeason(c echo.Context) error {
	var p *programmeSeasons.ProgrammeSeason
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("EditProgrammeSeason failed to get programme season: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.programmeSeasons.EditProgrammeSeason(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) DeleteProgrammeSeason(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	_, err = r.programmeSeasons.GetProgrammeSeasonById(id)
	if err != nil {
		err = fmt.Errorf("DeleteProgrammeSeason failed to get programme season: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.programmeSeasons.DeleteProgrammeSeason(id)
	if err != nil {
		err = fmt.Errorf("DeleteProgrammeSeason failed to delete programme season: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
