package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/COMTOP1/api/services/admin"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

type (
	Repo struct {
		admin        *admin.Store
		access       *utils.Accesser
		domainName   string
		signingToken string
	}

	Claim struct {
		jwt.StandardClaims
	}

	Token struct {
		Raw       string                 // The raw token.  Populated when you Parse a token
		Method    SigningMethod          // The signing method used or to be used
		Header    map[string]interface{} // The first segment of the token
		Claims    Claims                 // The second segment of the token
		Signature string                 // The third segment of the token.  Populated when you Parse a token
		Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
	}

	SigningMethod interface {
		Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
		Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
		Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
	}

	Claims interface {
		Valid() error
	}

	JWTToken struct {
		JWTToken string `json:"jwt_token"`
	}
)

func NewRepo(scope *gocb.Scope, access *utils.Accesser, domainName string, signingToken string) *Repo {
	return &Repo{
		admin:        admin.NewStore(scope),
		access:       access,
		domainName:   domainName,
		signingToken: signingToken,
	}
}

func (r *Repo) GetJWT(c echo.Context) error {
	timeNow := time.Now()
	expirationTime := timeNow.Add(5 * time.Minute).Unix()
	uuid1 := uuid.NewV4()
	claim := &jwt.StandardClaims{
		Audience:  "https://api." + r.domainName,
		ExpiresAt: expirationTime,
		Id:        uuid1.String(),
		IssuedAt:  timeNow.Unix(),
		Issuer:    "https://api." + r.domainName,
		NotBefore: timeNow.Unix(),
	}
	token := NewWithClaims(jwt.SigningMethodHS512, claim)
	tokenString, err := token.SignedString([]byte(r.signingToken))
	if err != nil {
		err = fmt.Errorf("GetJWT failed: %w", err)
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = json.NewEncoder(c.Response().Writer).Encode(JWTToken{JWTToken: tokenString})
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (r *Repo) GetSSOJWT(c echo.Context) error {
	timeNow := time.Now()
	uuid1 := uuid.NewV4()
	claim := &jwt.StandardClaims{
		Audience:  "https://sso." + r.domainName,
		Id:        uuid1.String(),
		IssuedAt:  timeNow.Unix(),
		Issuer:    "https://api." + r.domainName,
		NotBefore: timeNow.Unix(),
	}
	token := NewWithClaims(jwt.SigningMethodHS512, claim)
	tokenString, err := token.SignedString([]byte(r.signingToken))
	if err != nil {
		err = fmt.Errorf("GetJWT failed: %w", err)
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	err = json.NewEncoder(c.Response().Writer).Encode(JWTToken{JWTToken: tokenString})
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, utils.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func NewWithClaims(method SigningMethod, claims Claims) *Token {
	uuid1 := uuid.NewV4()
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"kid": uuid1.String(),
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}

func (t *Token) SignedString(key interface{}) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = t.SigningString(); err != nil {
		return "", err
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", err
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

func (t *Token) SigningString() (string, error) {
	var err error
	parts := make([]string, 2)
	for i := range parts {
		var jsonValue []byte
		if i == 0 {
			fmt.Println(t.Header)
			if jsonValue, err = json.Marshal(t.Header); err != nil {
				return "", err
			}
		} else {
			fmt.Println(t.Claims)
			if jsonValue, err = json.Marshal(t.Claims); err != nil {
				return "", err
			}
		}

		parts[i] = EncodeSegment(jsonValue)
	}
	return strings.Join(parts, "."), nil
}

func EncodeSegment(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}
