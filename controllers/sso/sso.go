package sso

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/couchbase/gocb/v2"
)

type Repos struct {
	ServiceURL string
}

func NewRepo(scope *gocb.Scope, controller controllers.Controller, serviceURL string) *Repos {
	return &Repos{
		ServiceURL: serviceURL,
	}
}
