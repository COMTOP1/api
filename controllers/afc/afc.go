package afc

import (
	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/controllers/afc/v1/affiliations"
	"github.com/COMTOP1/api/controllers/afc/v1/documents"
	"github.com/COMTOP1/api/controllers/afc/v1/images"
	"github.com/COMTOP1/api/controllers/afc/v1/news"
	"github.com/COMTOP1/api/controllers/afc/v1/players"
	"github.com/COMTOP1/api/controllers/afc/v1/programmeSeasons"
	"github.com/COMTOP1/api/controllers/afc/v1/programmes"
	"github.com/COMTOP1/api/controllers/afc/v1/sponsors"
	"github.com/COMTOP1/api/controllers/afc/v1/teams"
	"github.com/COMTOP1/api/controllers/afc/v1/users"
	"github.com/COMTOP1/api/controllers/afc/v1/whatsOn"
	"github.com/couchbase/gocb/v2"
)

type Repos struct {
	ServiceURL       string
	Affiliations     *affiliations.Repo
	Documents        *documents.Repo
	Images           *images.Repo
	News             *news.Repo
	Players          *players.Repo
	Programmes       *programmes.Repo
	ProgrammeSeasons *programmeSeasons.Repo
	Sponsors         *sponsors.Repo
	Teams            *teams.Repo
	Users            *users.Repo
	WhatsOn          *whatsOn.Repo
}

func NewRepos(scope *gocb.Scope, controller controllers.Controller, serviceURL string) *Repos {
	return &Repos{
		ServiceURL:       serviceURL,
		Affiliations:     affiliations.NewRepo(scope, controller),
		Documents:        documents.NewRepo(scope, controller),
		Images:           images.NewRepo(scope, controller),
		News:             news.NewRepo(scope, controller),
		Players:          players.NewRepo(scope, controller),
		Programmes:       programmes.NewRepo(scope, controller),
		ProgrammeSeasons: programmeSeasons.NewRepo(scope, controller),
		Sponsors:         sponsors.NewRepo(scope, controller),
		Teams:            teams.NewRepo(scope, controller),
		Users:            users.NewRepo(scope, controller),
		WhatsOn:          whatsOn.NewRepo(scope, controller),
	}
}
