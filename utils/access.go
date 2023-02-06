package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
)

type (
	Accesser struct {
		conf Config
	}
	Config struct {
		AccessCookieName string
		SigningKey       []byte
		DomainName       string
		Admin            Admin
	}
	Admin struct {
		AdminAccessCookieName, Key0, Key1, Key2, Key3, TOTP string
	}
	AdminRequest struct {
		Key0     string `json:"key_0"`
		Key1     string `json:"key_1"`
		Key2     string `json:"key_2"`
		Key3     string `json:"key_3"`
		TOTPCode string `json:"totp_code"`
	}
	AdminResponse struct {
		JWTToken string `json:"jwt_token"`
	}
	AFCAccessClaims struct {
		Id   uint64 `json:"id"`
		Role string `json:"role"`
		jwt.StandardClaims
	}
	// AccessClaims represents an identifiable JWT
	AccessClaims struct {
		UserID      int      `json:"id"`
		Role        string   `json:"role"`
		Permissions []string `json:"perms"`
		jwt.StandardClaims
	}
	// Permission represents the permissions that a user has
	Permission struct {
		Name string `json:"name"`
	}
)

var (
	ErrNoToken      = errors.New("token not found")
	ErrInvalidToken = errors.New("invalid token")
)

// NewAccesser allows the validation of JWT tokens both as
// headers and as cookies
func NewAccesser(conf Config) *Accesser {
	return &Accesser{
		conf: conf,
	}
}

// FindAdminToken will return the claims from an AFC access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) FindAdminToken(r *http.Request) bool {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.Admin.AdminAccessCookieName)
		if err != nil {
			return false
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return false
		}
		token = splitToken[1]
	}

	if token == "" {
		return false
	}
	return true
}

// GetAdminToken will return the claims from an admin access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetAdminToken(r *http.Request) (*jwt.Token, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.Admin.AdminAccessCookieName)
		if err != nil {
			if errors.As(http.ErrNoCookie, &err) {
				return nil, ErrNoToken
			}
			return nil, fmt.Errorf("failed to get cookie: %w", err)
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return nil, ErrInvalidToken
		}
		token = splitToken[1]
	}

	if token == "" {
		return nil, ErrNoToken
	}

	return a.getAdminClaims(token)
}

// GetAdminToken will return the claims from an admin access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetAdminTokenKIDAndClaims(r *http.Request) (string, *jwt.StandardClaims, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.Admin.AdminAccessCookieName)
		if err != nil {
			if errors.As(http.ErrNoCookie, &err) {
				return "", nil, ErrNoToken
			}
			return "", nil, fmt.Errorf("failed to get cookie: %w", err)
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return "", nil, ErrInvalidToken
		}
		token = splitToken[1]
	}

	if token == "" {
		return "", nil, ErrNoToken
	}

	return a.getAdminClaimsKIDAndClaims(token)
}

// GetAFCToken will return the claims from an AFC access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetAFCToken(r *http.Request) (*AFCAccessClaims, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.AccessCookieName)
		if err != nil {
			if errors.As(http.ErrNoCookie, &err) {
				return nil, ErrNoToken
			}
			return nil, fmt.Errorf("failed to get cookie: %w", err)
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return nil, ErrInvalidToken
		}
		token = splitToken[1]
	}

	if token == "" {
		return nil, ErrNoToken
	}
	return a.getAFCClaims(token)
}

// GetToken will return the claims from an access token JWT
//
// First will check the Authorization header, if unset will
// check the access cookie
func (a *Accesser) GetToken(r *http.Request) (*AccessClaims, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		cookie, err := r.Cookie(a.conf.AccessCookieName)
		if err != nil {
			if errors.As(http.ErrNoCookie, &err) {
				return nil, ErrNoToken
			}
			return nil, fmt.Errorf("failed to get cookie: %w", err)
		}
		token = cookie.Value
	} else {
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			return nil, ErrInvalidToken
		}
		token = splitToken[1]
	}

	if token == "" {
		return nil, ErrNoToken
	}
	return a.getClaims(token)
}

func (a *Accesser) getAFCClaims(token string) (*AFCAccessClaims, error) {
	claims := &AFCAccessClaims{}

	jwt1, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if jwt1.Valid && strings.Contains(claims.Issuer, "https://sso."+a.conf.DomainName) && claims.Audience == "https://afcaldermaston.co.uk" && claims.IssuedAt == claims.NotBefore && claims.ExpiresAt > time.Now().Unix() {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (a *Accesser) getClaims(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	jwt1, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if jwt1.Valid && claims.Issuer == "https://sso."+a.conf.DomainName && strings.Contains(claims.Audience, a.conf.DomainName) && claims.IssuedAt == claims.NotBefore && claims.ExpiresAt > time.Now().Unix() {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (a *Accesser) getAdminClaims(token string) (*jwt.Token, error) {
	claims := &jwt.StandardClaims{}

	jwt1, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if jwt1.Valid && claims.Issuer == "https://sso."+a.conf.DomainName && strings.Contains(claims.Audience, a.conf.DomainName) && claims.IssuedAt == claims.NotBefore && claims.ExpiresAt > time.Now().Unix() {
		return jwt1, nil
	}
	return nil, ErrInvalidToken
}

func (a *Accesser) getAdminClaimsKIDAndClaims(token string) (string, *jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}

	jwt1, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.conf.SigningKey, nil
	})
	if err != nil {
		return "", nil, ErrInvalidToken
	}
	if jwt1.Valid && claims.Issuer == "https://sso."+a.conf.DomainName && strings.Contains(claims.Audience, a.conf.DomainName) && claims.IssuedAt == claims.NotBefore && claims.ExpiresAt > time.Now().Unix() {
		return fmt.Sprintf("%v", jwt1.Header["kid"]), claims, nil
	}
	return "", nil, ErrInvalidToken
}

// AdminInitAuthMiddleware checks a HTTP request for a valid token either in the header or cookie
func (a *Accesser) AdminInitAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if a.FindAdminToken(c.Request()) {
			return echo.NewHTTPError(http.StatusUnauthorized, "token already exists")
			//return &echo.HTTPError{
			//	Code:    http.StatusUnauthorized,
			//	Message: "token already exists",
			//}
		}
		var adminRequest *AdminRequest
		err := json.NewDecoder(c.Request().Body).Decode(&adminRequest)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}
		valid := totp.Validate(adminRequest.TOTPCode, a.conf.Admin.TOTP)
		if !valid {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "totp invalid",
				Internal: fmt.Errorf("inputed totp code is invalid"),
			}
		}
		if adminRequest.Key0 != a.conf.Admin.Key0 || adminRequest.Key1 != a.conf.Admin.Key1 || adminRequest.Key2 != a.conf.Admin.Key2 || adminRequest.Key3 != a.conf.Admin.Key3 {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "key invalid",
				Internal: fmt.Errorf("inputed admin key is invalid"),
			}
		}
		return next(c)
	}
}

// AdminAuthMiddleware checks a HTTP request for a valid token either in the header or cookie
func (a *Accesser) AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.GetAdminToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}
		return next(c)
	}
}

// AFCAuthMiddleware checks a HTTP request for a valid token either in the header or cookie
func (a *Accesser) AFCAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.GetAFCToken(c.Request())
		if err != nil {
			//return echo.NewHTTPError(http.StatusUnauthorized, err)
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}
		return next(c)
	}
}

func (a *Accesser) AFCAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := a.GetAFCToken(c.Request())
		if err != nil {
			//return echo.NewHTTPError(http.StatusUnauthorized, err)
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}
		if token.Role == "ClubSecretary" || token.Role == "Chairperson" || token.Role == "Webmaster" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  fmt.Errorf("not authenticated"),
			Internal: fmt.Errorf("not authenticated"),
		}
	}
}

// AuthMiddleware checks a HTTP request for a valid token either in the header or cookie
func (a *Accesser) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := a.GetToken(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}
		return next(c)
	}
}

// NilMiddleware checks a HTTP request for a valid token either in the header or cookie
func (a *Accesser) NilMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
