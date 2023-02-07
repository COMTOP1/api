package teams

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/teams"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"unicode"
)

type Repo struct {
	teams      *teams.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		teams:      teams.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) GetTeamById(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("GetTeamById failed to get id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	w, err := r.teams.GetTeamById(id)
	if err != nil {
		err = fmt.Errorf("GetTeamById failed to get whatson: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) ListAllTeams(c echo.Context) error {
	t, err := r.teams.ListAllTeams()
	if err != nil {
		err = fmt.Errorf("ListAllTeams failed to get all teams: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (r *Repo) ListActiveTeams(c echo.Context) error {
	t, err := r.teams.ListActiveTeams()
	if err != nil {
		err = fmt.Errorf("ListActiveTeams failed to get active teams: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (r *Repo) AddTeam(c echo.Context) error {
	var t *teams.Team
	err := c.Bind(&t)
	if err != nil {
		err = fmt.Errorf("AddTeam failed to bind team: %t", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.teams.AddTeam(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (r *Repo) EditTeam(c echo.Context) error {
	var t *teams.Team
	err := c.Bind(&t)
	if err != nil {
		err = fmt.Errorf("EditTeam failed to get team: %t", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.teams.EditTeam(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (r *Repo) DeleteTeam(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	_, err = r.teams.GetTeamById(id)
	if err != nil {
		err = fmt.Errorf("DeleteTeam failed to get user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.teams.DeleteTeam(id)
	if err != nil {
		err = fmt.Errorf("DeleteTeam failed to delete user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
