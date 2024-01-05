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
	"github.com/couchbase/gocb/v2"
)

type (
	AffiliationsRepo interface {
		affiliations.Repo
	}

	DocumentsRepo interface {
		documents.Repo
	}

	ImagesRepo interface {
		images.Repo
	}

	NewsRepo interface {
		news.Repo
	}

	PlayersRepo interface {
		players.Repo
	}

	ProgrammesRepo interface {
		programmes.Repo
	}

	SponsorsRepo interface {
		sponsors.Repo
	}

	TeamsRepo interface {
		teams.Repo
	}

	UsersRepo interface {
		users.Repo
	}

	WhatsOnRepo interface {
		whatsOn.Repo
	}

	//nolint:unused
	Store struct {
		//nolint:unused
		scope *gocb.Scope
	}
)
