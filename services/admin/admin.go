package admin

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/golang-jwt/jwt"
	"strings"
)

type (
	AdminRepo interface {
		VerifySSOJWT(claims *jwt.StandardClaims, kid, userAgent string) error
		AddSSOJWT(claims *jwt.StandardClaims, kid, userAgent string) error
	}

	// Store contains our dependency
	Store struct {
		scope *gocb.Scope
	}

	JWTToken struct {
		KID       string `json:"kid"`
		JTI       string `json:"jti"`
		IAT       int64  `json:"iat"`
		NBF       int64  `json:"nbf"`
		EXP       int64  `json:"exp"`
		UserAgent string `json:"userAgent"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (s *Store) VerifySSOJWT(claims *jwt.StandardClaims, kid, userAgent string) error {
	get, err := s.scope.Collection("sso").Get("sso:"+kid, &gocb.GetOptions{})
	if err != nil {
		return err
	}
	var token JWTToken
	err = get.Content(&token)
	if err != nil {
		return fmt.Errorf("failed to get jwt: %w", err)
	}
	if !(token.KID == kid && token.JTI == claims.Id && token.IAT == claims.IssuedAt && token.NBF == claims.NotBefore && token.UserAgent == userAgent) {
		return fmt.Errorf("failed to validate jwt")
	}
	return nil
}

func (s *Store) AddSSOJWT(claims *jwt.StandardClaims, kid, userAgent string) error {
	query, err := s.scope.Query("SELECT `kid` FROM sso", &gocb.QueryOptions{})
	if err != nil {
		return err
	}
	for query.Next() {
		var result JWTToken
		err = query.Row(&result)
		if err != nil {
			if !strings.Contains(err.Error(), "document not found") {
				return fmt.Errorf("failed to get sso in add sso jwt: %w", err)
			}
		}
		_, err = s.scope.Collection("sso").Remove("sso:"+result.KID, nil)
		if err != nil {
			return err
		}
	}
	jwtToken := JWTToken{
		KID:       kid,
		JTI:       claims.Id,
		IAT:       claims.IssuedAt,
		NBF:       claims.NotBefore,
		EXP:       claims.ExpiresAt,
		UserAgent: userAgent,
	}
	_, err = s.scope.Collection("sso").Insert("sso:"+kid, jwtToken, nil)
	if err != nil {
		return fmt.Errorf("failed to add whatsOn: %w", err)
	}
	return nil
}
