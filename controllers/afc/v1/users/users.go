package users

import (
	"fmt"
	"github.com/COMTOP1/api/services/afc/users"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Repo stores our dependencies
type Repo struct {
	users  *users.Store
	access *utils.Accesser
}

// NewRepo creates our data store
func NewRepo(scope *gocb.Scope, access *utils.Accesser) *Repo {
	return &Repo{
		users:  users.NewStore(scope),
		access: access,
	}
}

// UserByEmail finds a user by email
// @Summary Get a user by email
// @Description Get a basic user object by email.
// @ID get-user-email
// @Tags user-email
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} users.User
// @Router /ea231a602d352b2bcc5a2acca6022575/v1/internal/user/{email} [get]
func (r *Repo) UserByEmail(c echo.Context) error {
	email := c.Param("email")
	p, err := r.users.GetUser(email)
    if err != nil {
        err = fmt.Errorf("UserByEmail failed: %w", err)
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
    }
	return c.JSON(http.StatusOK, p)
}

// UserByEmailFull finds a user by email returning all info
// @Summary Get a full user by email
// @Description Get a complete user object by email.
// @ID get-user-email-full
// @Tags user-email-full
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} users.User
// @Router /ea231a602d352b2bcc5a2acca6022575/v1/internal/user/{email}/full [get]
func (r *Repo) UserByEmailFull(c echo.Context) error {
	email := c.Param("email")
	p, err := r.users.GetUserFull(email)
	if err != nil {
		err = fmt.Errorf("UserByEmailFull failed to get user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

// UserByToken finds a user by their JWT token
// @Summary Get a user by token
// @Description Get a basic user object by JWT token generated by web-auth.
// @ID get-user-token
// @Tags user-token
// @Produce json
// @Success 200 {object} users.User
// @Router /ea231a602d352b2bcc5a2acca6022575/v1/internal/user [get]
func (r *Repo) UserByToken(c echo.Context) error {
	claims, err := r.access.GetAFCToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByToken failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, utils.Error{Error: err.Error()})
	}
	p, err := r.users.GetUser(claims.Email)
	if err != nil {
		err = fmt.Errorf("UserByToken failed getting user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

// UserByTokenFull finds a user by their JWT token returning all info
// @Summary Get a full user by token
// @Description Get a complete user object by JWT token generated by web-auth.
// @ID get-user-token-full
// @Tags user-token-full
// @Produce json
// @Success 200 {object} user.UserFull
// @Router /ea231a602d352b2bcc5a2acca6022575/v1/internal/user/full [get]
func (r *Repo) UserByTokenFull(c echo.Context) error {
	claims, err := r.access.GetAFCToken(c.Request())
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	p, err := r.users.GetUserFull(claims.Email)
	if err != nil {
		err = fmt.Errorf("UserByTokenFull failed getting user: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

// ListAllUsers handles listing all users
//
// @Summary List all users
// @ID get-users-all
// @Tags users-all
// @Produce json
// @Success 200 {array} users.User
// @Router /ea231a602d352b2bcc5a2acca6022575/v1/internal/users [get]
func (r *Repo) ListAllUsers(c echo.Context) error {
	u, err := r.users.ListAllUsers()
	if err != nil {
		err = fmt.Errorf("ListAllUsers failed to get all users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, u)
}

// ListAllContactUsers handles listing all contact users
//
// @Summary List all contact users
// @ID get-users-contact-all
// @Tags users-all-contact
// @Produce json
// @Success 200 {array} users.User
// @Router /ea231a602d352b2bcc5a2acca6022575/v1/public/contacts [get]
func (r *Repo) ListAllContactUsers(c echo.Context) error {
	u, err := r.users.ListContactUsers()
	if err != nil {
		err = fmt.Errorf("ListAllContactUsers failed to get all contact users: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, u)
}

func (r *Repo) AddUser(c echo.Context) error {
	var u *users.UserFull
	err := c.Bind(&u)
	if err != nil {
		err = fmt.Errorf("AddUser failed to get user: %w", err)
		return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = r.users.AddUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, u)
}

func (r *Repo) EditUser(c echo.Context) error {
    var u *users.UserFull
    err := c.Bind(&u)
    if err != nil {
        err = fmt.Errorf("EditUser failed to get user: %w", err)
        return c.JSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
    }
    err = r.users.EditUser(u)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
    }
    return c.JSON(http.StatusOK, u)
}

func (r *Repo) DeleteUser(c echo.Context) error {
    email := c.Param("email")
    _, err := r.users.GetUserFull(email)
    if err != nil {
        err = fmt.Errorf("DeleteUser failed to get user: %w", err)
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
    }
    err = r.users.DeleteUser(email)
    if err != nil {
        err = fmt.Errorf("DeleteUser failed to delete user: %w", err)
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
    }
    return c.NoContent(http.StatusOK)
}
