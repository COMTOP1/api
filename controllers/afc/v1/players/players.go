package players

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/players"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"unicode"
)

type Repo struct {
	players    *players.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		players:    players.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) GetPlayerById(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("GetPlayerById failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.players.GetPlayerById(id)
	if err != nil {
		err = fmt.Errorf("GetPlayerById failed to get player: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllPlayers(c echo.Context) error {
	p, err := r.players.ListAllPlayers()
	if err != nil {
		err = fmt.Errorf("ListAllPlayers failed to get all player: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllPlayersByTeamId(c echo.Context) error {
	temp := c.Param("teamId")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	teamId, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("ListAllPlayersByTeamId failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.players.ListAllPlayersByTeamId(teamId)
	if err != nil {
		err = fmt.Errorf("ListAllPlayersByTeamId failed to get all player: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) AddPlayer(c echo.Context) error {
	var p *players.Player
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("AddPlayer failed to bind player: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.players.AddPlayer(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) EditPlayer(c echo.Context) error {
	var p *players.Player
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("EditPlayer failed to bind player: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.players.EditPlayer(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) DeletePlayer(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	_, err = r.players.GetPlayerById(id)
	if err != nil {
		err = fmt.Errorf("DeletePlayer failed to get player: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.players.DeletePlayer(id)
	if err != nil {
		err = fmt.Errorf("DeletePlayer failed to delete player: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
