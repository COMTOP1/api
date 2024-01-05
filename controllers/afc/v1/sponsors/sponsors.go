package sponsors

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/sponsors"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"unicode"
)

type Repo struct {
	sponsors   *sponsors.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		sponsors:   sponsors.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) GetSponsorById(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("GetSponsorById failed to get id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	w, err := r.sponsors.GetSponsorById(id)
	if err != nil {
		err = fmt.Errorf("GetSponsorById failed to get sponsor: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, w)
}

func (r *Repo) ListALlSponsors(c echo.Context) error {
	s, err := r.sponsors.ListAllSponsors()
	if err != nil {
		err = fmt.Errorf("ListAllSponsor failed to get all sponsors: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func (r *Repo) ListALlSponsorsMinimal(c echo.Context) error {
	s, err := r.sponsors.ListAllSponsorsMinimal()
	if err != nil {
		err = fmt.Errorf("ListAllSponsorMinimal failed to get all sponsors minimal: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func (r *Repo) ListAllSponsorsByTeamId(c echo.Context) error {
	teamId := c.Param("teamId")
	s, err := r.sponsors.ListAllSponsorsByTeamId(teamId)
	if err != nil {
		err = fmt.Errorf("ListAllSponsorsByTeamId failed to get all sponsors by team id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func (r *Repo) AddSponsor(c echo.Context) error {
	var s *sponsors.Sponsor
	err := c.Bind(&s)
	if err != nil {
		err = fmt.Errorf("AddSponsor failed to bind sponsor: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.sponsors.AddSponsor(s)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func (r *Repo) EditSponsor(c echo.Context) error {
	var s *sponsors.Sponsor
	err := c.Bind(&s)
	if err != nil {
		err = fmt.Errorf("EditSponsor failed to get sponsor: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.sponsors.EditSponsor(s)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func (r *Repo) DeleteSponsor(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("DeleteSponsor failed to get id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	_, err = r.sponsors.GetSponsorById(id)
	if err != nil {
		err = fmt.Errorf("DeleteSponsor failed to get sponsor: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.sponsors.DeleteSponsor(id)
	if err != nil {
		err = fmt.Errorf("DeleteSponsor failed to delete sponsor: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
