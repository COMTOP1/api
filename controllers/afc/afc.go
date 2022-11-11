package afc

import (
	"github.com/COMTOP1/api/controllers/afc/v1/affiliations"
	"github.com/COMTOP1/api/controllers/afc/v1/documents"
	"github.com/COMTOP1/api/controllers/afc/v1/images"
	"github.com/COMTOP1/api/controllers/afc/v1/news"
	"github.com/COMTOP1/api/controllers/afc/v1/players"
	"github.com/COMTOP1/api/controllers/afc/v1/programmes"
	"github.com/COMTOP1/api/controllers/afc/v1/sponsors"
	"github.com/COMTOP1/api/controllers/afc/v1/teams"
	"github.com/COMTOP1/api/controllers/afc/v1/users"
	"github.com/COMTOP1/api/controllers/afc/v1/whatsOn"
	"github.com/COMTOP1/api/utils"
	"github.com/couchbase/gocb/v2"
)

type Repos struct {
	access       *utils.Accesser
	Affiliations *affiliations.Repo
	Documents    *documents.Repo
	Images       *images.Repo
	News         *news.Repo
	Players      *players.Repo
	Programmes   *programmes.Repo
	Sponsors     *sponsors.Repo
	Teams        *teams.Repo
	Users        *users.Repo
	WhatsOn      *whatsOn.Repo
}

func NewRepos(scope *gocb.Scope, access *utils.Accesser) *Repos {
	return &Repos{
		access:       access,
		Affiliations: affiliations.NewRepo(scope, access),
		Documents:    documents.NewRepo(scope, access),
		Images:       images.NewRepo(scope, access),
		News:         news.NewRepo(scope, access),
		Players:      players.NewRepo(scope, access),
		Programmes:   programmes.NewRepo(scope, access),
		Sponsors:     sponsors.NewRepo(scope, access),
		Teams:        teams.NewRepo(scope, access),
		Users:        users.NewRepo(scope, access),
		WhatsOn:      whatsOn.NewRepo(scope, access),
	}
}
