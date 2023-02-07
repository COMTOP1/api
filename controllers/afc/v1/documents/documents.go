package documents

import (
	"fmt"
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/services/afc/documents"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"unicode"
)

type Repo struct {
	documents  *documents.Store
	controller controllers.Controller
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller) *Repo {
	return &Repo{
		documents:  documents.NewStore(scope),
		controller: controller,
	}
}

func (r *Repo) GetDocumentById(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	if err != nil {
		err = fmt.Errorf("GetDocumentById failed to get id: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.documents.GetDocumentById(id)
	if err != nil {
		err = fmt.Errorf("GetDocumentById failed to get document: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (r *Repo) ListAllDocuments(c echo.Context) error {
	d, err := r.documents.ListAllDocuments()
	if err != nil {
		err = fmt.Errorf("ListAllDocuments failed to get all document: %p", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (r *Repo) AddDocument(c echo.Context) error {
	var d *documents.Document
	err := c.Bind(&d)
	if err != nil {
		err = fmt.Errorf("AddDocument failed to bind document: %p", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.documents.AddDocument(d)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (r *Repo) DeleteDocument(c echo.Context) error {
	temp := c.Param("id")
	temp1 := []rune(temp)
	for _, r2 := range temp1 {
		if !unicode.IsNumber(r2) {
			return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: "id expects a positive number, the provided is not a positive number"})
		}
	}
	id, err := strconv.ParseUint(temp, 10, 64)
	_, err = r.documents.GetDocumentById(id)
	if err != nil {
		err = fmt.Errorf("DeleteDocument failed to get document: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.documents.DeleteDocument(id)
	if err != nil {
		err = fmt.Errorf("DeleteDocument failed to delete document: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
