package affiliations

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/affiliations"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"unicode"
)

// Repo stores our dependencies
type Repo struct {
	affiliations *affiliations.Store
	controller   controllers.Controller
}

// NewRepo creates our data store
func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		affiliations: affiliations.NewStore(scope),
		controller:   controller,
	}
}

func (r *Repo) GetAffiliationById(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("GetAffiliationById failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.affiliations.GetAffiliationById(id)
	if err != nil {
		err = fmt.Errorf("GetAffiliationById failed to get affiliation: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllAffiliations(c echo.Context) error {
	a, err := r.affiliations.ListAllAffiliations()
	if err != nil {
		err = fmt.Errorf("ListAllAffiliations failed to get all affiliation: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, a)
}

func (r *Repo) AddAffiliation(c echo.Context) error {
	var a *affiliations.Affiliation
	err := c.Bind(&a)
	if err != nil {
		err = fmt.Errorf("AddAffiliation failed to bind affiliation: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.affiliations.AddAffiliation(a)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, a)
}

func (r *Repo) DeleteAffiliation(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("DeleteAffiliation failed to get id: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	_, err = r.affiliations.GetAffiliationById(id)
	if err != nil {
		err = fmt.Errorf("DeleteAffiliation failed to get affiliation: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.affiliations.DeleteAffiliation(id)
	if err != nil {
		err = fmt.Errorf("DeleteAffiliation failed to delete affiliation: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
