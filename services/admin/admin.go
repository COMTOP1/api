package admin

import (
    "fmt"
    "github.com/couchbase/gocb/v2"
    "github.com/golang-jwt/jwt"
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
		KID string `json:"kid"`
        JTI string `json:"jti"`
        IAT int64 `json:"iat"`
        NBF int64 `json:"nbf"`
        UserAgent string `json:"userAgent"`
	}
)

func NewStore(scope *gocb.Scope) *Store {
	return &Store{scope: scope}
}

func (s *Store) VerifySSOJWT(claims *jwt.StandardClaims, kid, userAgent string) error {
    query, err := s.scope.Query("SELECT `kid` FROM sso WHERE `kid` = $1", &gocb.QueryOptions{
        Adhoc: false,
        PositionalParameters: []interface{}{kid},
        })
    if err != nil {
        return err
    }
    for query.Next() {
        var result JWTToken
        err = query.Row(&result)
        if err != nil {
            return fmt.Errorf("failed to get jwt: %w", err)
        }
        if !(result.KID == kid && result.JTI == claims.Id && result.IAT == claims.IssuedAt && result.NBF == claims.NotBefore && result.UserAgent == userAgent) {
            return fmt.Errorf("failed to validate jwt")
        }
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
            return fmt.Errorf("failed to get jwt: %w", err)
        }
        mut, err := s.scope.Collection("sso").Remove("sso:"+kid, nil)
        fmt.Println(mut)
        if err != nil {
            return err
        }
    }
    jwtToken := JWTToken{
        KID: kid,
        JTI: claims.Id,
        IAT: claims.IssuedAt,
        NBF: claims.NotBefore,
        UserAgent: userAgent,
    }
    mut, err := s.scope.Collection("sso").Insert("sso:"+kid, jwtToken, nil)
    fmt.Println(mut)
    if err != nil {
        return fmt.Errorf("failed to add whatsOn: %w", err)
    }
    return nil
}
