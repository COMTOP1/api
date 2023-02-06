package routes

import (
	"fmt"
	"github.com/COMTOP1/api/controllers/afc"
	"net/http"
	"time"

	adminPackage "github.com/COMTOP1/api/controllers/admin/v1/admin"
	"github.com/COMTOP1/api/middleware"
	"github.com/COMTOP1/api/utils"
	"github.com/labstack/echo/v4"
)

type (
	Router struct {
		port    string
		version string
		commit  string
		router  *echo.Echo
		access  *utils.Accesser
		mailer  *utils.Mailer
		afc     *afc.Repos
		admin   *adminPackage.Repo
	}

	NewRouter struct {
		Port       string
		Version    string
		Commit     string
		DomainName string
		Debug      bool
		Access     *utils.Accesser
		Mailer     *utils.Mailer
		AFC        *afc.Repos
		Admin      *adminPackage.Repo
	}
)

func New(conf *NewRouter) *Router {
	r := &Router{
		port:    conf.Port,
		version: conf.Version,
		commit:  conf.Commit,
		router:  echo.New(),
		access:  conf.Access,
		mailer:  conf.Mailer,
		afc:     conf.AFC,
		admin:   conf.Admin,
	}
	r.router.HideBanner = true

	r.router.Debug = conf.Debug

	middleware.New(r.router, conf.DomainName)

	r.loadRoutes()

	return r
}

func (r *Router) Start() error {
	r.router.Logger.Error(r.router.Start(r.port))
	return fmt.Errorf("failed to start router on port %s", r.port)
}

// loadRoutes initialise routes
// @title BSWDI API
// @description The backend powering most things
// @contact.name API Support
// @contact.url https://github.com/COMTOP1/api
// @contact.email api@bswdi.co.uk
func (r *Router) loadRoutes() {
	r.router.RouteNotFound("/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, utils.Error{Error: "Not found"})
	})

	r.router.GET("/", func(c echo.Context) error {
		text := fmt.Sprintf(`                               -+yysssssosyyyyy+-
                             .+ys/.'''''-oyyo:/yy+.
                          '.+ss+.     -+yyo-'  .+ss+.
                         .+ss+.     .+syo:'      .+ss+.
                       ./ss+.'    .+syo:'    '    '.+ss+.
                     ./ss+-'    .+syo:'    '-o/'    '-+ss/.
                   ./ss+-'    ./sys:'    '-oyyys/.    '-oys/.
                 ./sy+-'    ./sys:.    '-oys+-/sys/'    '-oys/.
               ':syo-'    ./sys/.    '-+yy+-'  ./sys:'    ':oys/'
             ':oyo:'    'v:sys/.    '-+yssy+.'    ./sys:'    ':oyo:'
           ':oyo:'    ':sys/.    '.+yy+--+ss/.     .+syo:'    ':oyo:'
         ':oys:'    ':sys+.     .+yyo-    .+ys/'     .+syo:'    ':syo:'
       ':oys:'    ':oyy+.     .+yyo-        .+ys/'     .+yyo:'    ':syo-'
     '-oys:'    ':oyy+.     .+syo-'           .oys/'     .+yyo-'    ':syo-'
   '-oys/'    '-oyy+.     ./yyo-'              '-oys:'     -+yyo-'    '/syo-'
 '-oys/'    '-oyy+-     ./syo-'                  '-oys:'     -+yyo-'    '/syo-'
-oys/.    '-oyy+-     '/syo:'                      ':oys:'    '-oyyo-'    ./sy+-
yys:'     -syy/'    '/syo:'                          ':oyo:'    '-oys+-     .+yy
yyys+.     ./syo-''/syo:'                              ':oyo:'    '-oys+-    .yy
yy/oys+.     .+ssosys:'                                  ':oyo:'    ':oys+.  .yy
yy..:sys/.    '.+sys:'                                    ':syyo-'    ':oys/..yy
yy.  ./sys/.    '-+ss+-'                                '.+sy+/syo-'    ':oys+yy
yy.    ./sys/.     -+ys+.'                             .+sy+-' ./syo-'    ':syyy
syo-'    ./sys/.     -+ys+.'                        '.+syo-'    .+yyo.     ':sys
./sy+.'    ./sys/'    '-oys+.                      ./syo-'    ':sys/.    '.+ys/.
  ./ss+.     .+sys:'    '-oys/.                  ./syo-'    ':sys/.    '.+ss/.
    .+ss+.     .+syo:'    '-oys/.              '/sys:'    ':oys/.     .+ss+.
      .+yy+.     .+yyo:'    '-oys/'          ':sys:'    ':oyy/.     .+sy+.
        .+yy/.     .+yyo:'    ':oys/'      ':sys:'    '-oyy/.     ./sy+.
          .+ys/.     -+yyo-'    ':oys/'  ':sys/'    '-oyy+.     ./sy+.'
           '-+ys/.    '-+yyo-'    ':sys/:oys/'    '-oys+.     ./sy+-'
             '-+ys/'    '-oyyo-'    '/yyys/.    '-oys+.'    ./sy+-'
               '-oys/'    '-oyy+-' ':oys/.    '-+ys+-'    '/sso-'
                 '-oyo:'    ':oyy+:oys/.    '-+yyo-'    ':oyo-'
                   ':oyo:'    ':oyys+.    '.+yy+-'    ':oyo-'
                     ':oyo:'    ':/.    '.+syo-'    ':oyo:'
                       ':oyo-'         .+syo:'    '-oyo:'
                         ':oyo-'     ./syo:'    '-oyo:'
                           '/syo-' ./sys:'    '-oys:'    BSWDI API
                             '/sy+/sys/.'''''-+ys/'      Version: %s
                               '/yyyyyyyyyyyyys/'        Commit ID: %s                                
`, r.version, r.commit)
		return c.String(http.StatusOK, text)
	})

	r.router.GET("/ping", func(c echo.Context) error {
		resp := map[string]time.Time{"pong": time.Now()}
		return c.JSON(http.StatusOK, resp)
	})

	r.loadAdminRoutes()

	r.loadAFCRoutes()

	r.loadBSWDIRoutes()
}

func (r *Router) loadAdminRoutes() {
	admin := r.router.Group("/" + r.admin.GetAdminURL())
	{
		getJWT := admin.Group("/get")
		{
			if !r.router.Debug {
				getJWT.Use(r.access.AdminInitAuthMiddleware)
			}
			getJWT.GET("/jwt", r.admin.GetJWT)
		}
		sso := admin.Group("/sso")
		{
			if !r.router.Debug {
				sso.Use(r.access.AdminAuthMiddleware)
			}
			sso.GET("/jwt", r.admin.GetSSOJWT)
			sso.GET("/verify", r.admin.VerifySSO)
		}
	}
}

func (r *Router) loadAFCRoutes() {
	afcGroup := r.router.Group("/ea231a602d352b2bcc5a2acca6022575") // afc_group --> MD5
	{
		afcV1 := afcGroup.Group("/v1")
		{
			internal := afcV1.Group("/internal")
			{
				if !r.router.Debug {
					internal.Use(r.access.AFCAuthMiddleware)
				}
				affiliation := internal.Group("/affiliation")
				{
					affiliation.PUT("", r.afc.Affiliations.AddAffiliation)
					affiliation.DELETE("/:id", r.afc.Affiliations.DeleteAffiliation)
				}
				document := internal.Group("/document")
				{
					document.PUT("", r.afc.Documents.AddDocument)
					document.DELETE("/:id", r.afc.Documents.DeleteDocument)
				}
				image := internal.Group("/image")
				{
					image.PUT("", r.afc.Images.AddImage)
					image.DELETE("/:id", r.afc.Images.DeleteImage)
				}
				news := internal.Group("/news")
				{
					news.PUT("", r.afc.News.AddNews)
					news.PATCH("", r.afc.News.EditNews)
					news.DELETE("/:id", r.afc.News.DeleteNews)
				}
				player := internal.Group("/player")
				{
					player.PUT("", r.afc.Players.AddPlayer)
					player.PATCH("", r.afc.Players.EditPlayer)
					player.DELETE("/:id", r.afc.Players.DeletePlayer)
				}
				internal.GET("/players", r.afc.Players.ListAllPlayers)
				programme := internal.Group("/programme")
				{
					programme.PUT("", r.afc.Programmes.AddProgramme)
					programme.DELETE("/:id", r.afc.Programmes.DeleteProgramme)
				}
				programmeSeason := internal.Group("/programmeSeason")
				{
					programmeSeason.PUT("", r.afc.ProgrammeSeasons.AddProgrammeSeason)
					programmeSeason.PATCH("", r.afc.ProgrammeSeasons.EditProgrammeSeason)
					programmeSeason.DELETE("/:id", r.afc.ProgrammeSeasons.DeleteProgrammeSeason)
				}
				sponsor := internal.Group("/sponsor")
				{
					sponsor.PUT("", r.afc.Sponsors.AddSponsor)
					sponsor.PATCH("", r.afc.Sponsors.EditSponsor)
					sponsor.DELETE("/:id", r.afc.Sponsors.DeleteSponsor)
				}
				internal.GET("/teams", r.afc.Teams.ListAllTeams)
				team := internal.Group("/team")
				{
					team.PUT("", r.afc.Teams.AddTeam)
					team.PATCH("", r.afc.Teams.EditTeam)
					team.DELETE("/:id", r.afc.Teams.DeleteTeam)
				}
				user := internal.Group("/user")
				{
					admin := user.Group("/admin")
					{
						if !r.router.Debug {
							admin.Use(r.access.AFCAdminMiddleware)
						}
						admin.PUT("", r.afc.Users.AddUser)
						admin.PATCH("", r.afc.Users.EditUser)
						admin.DELETE("/email/:email", r.afc.Users.DeleteUserFromEmail)
						admin.DELETE("/id/:id", r.afc.Users.DeleteUserFromId)
					}
					email := user.Group("/email")
					{
						email.GET("/:email/full", r.afc.Users.GetUserByEmailFull)
						email.GET("/:email", r.afc.Users.GetUserByEmail)
					}
					id := user.Group("/id")
					{
						id.GET("/:id/full", r.afc.Users.GetUserByIdFull)
						id.GET("/:id", r.afc.Users.GetUserById)
					}
					user.GET("/all", r.afc.Users.ListAllUsers)
					user.GET("/full", r.afc.Users.GetUserByTokenFull)
					user.GET("", r.afc.Users.GetUserByToken)
				}
				whatsOn := internal.Group("/whatsOn")
				{
					_ = whatsOn
					whatsOn.PUT("", r.afc.WhatsOn.AddWhatsOn)
					whatsOn.PATCH("", r.afc.WhatsOn.EditWhatsOn)
					whatsOn.DELETE("/:id", r.afc.WhatsOn.DeleteWhatsOn)
				}
			}
			public := afcV1.Group("/public")
			{
				public.GET("/affiliation", r.afc.Affiliations.GetAffiliationById)
				public.GET("/affiliations", r.afc.Affiliations.ListAllAffiliations)
				public.GET("/contacts", r.afc.Users.ListContactUsers)
				public.GET("/document/:id", r.afc.Documents.GetDocumentById)
				public.GET("/documents", r.afc.Documents.ListAllDocuments)
				public.GET("/image/:id", r.afc.Images.GetImageById)
				public.GET("/images", r.afc.Images.ListAllImages)
				news := public.Group("/news")
				{
					news.GET("", r.afc.News.ListAllNews)
					news.GET("/latest", r.afc.News.GetNewsLatest)
					news.GET("/:id", r.afc.News.GetNewsById)
				}
				public.GET("/player/:id", r.afc.Players.GetPlayerById)
				public.GET("/players/:teamID", r.afc.Players.ListAllPlayersByTeamId)
				public.GET("/programme/:id", r.afc.Programmes.GetProgrammeById)
				public.GET("/programmes", r.afc.Programmes.ListAllProgrammes)
				public.GET("/programmeSeasons", r.afc.ProgrammeSeasons.ListAllProgrammeSeasons)
				public.GET("/programmeSeason/:id", r.afc.ProgrammeSeasons.GetProgrammeSeasonById)
				public.GET("/sponsor/:id", r.afc.Sponsors.GetSponsorById)
				sponsors := public.Group("/sponsors")
				{
					sponsors.GET("", r.afc.Sponsors.ListALlSponsors)
					sponsors.GET("/minimal", r.afc.Sponsors.ListALlSponsorsMinimal)
					sponsors.GET("/:teamID", r.afc.Sponsors.ListAllSponsorsByTeamId)
				}
				public.GET("/teams", r.afc.Teams.ListActiveTeams)
				team := public.Group("/team")
				{
					team.GET("/:id", r.afc.Teams.GetTeamById)
					team.GET("/managers/:teamId", r.afc.Users.ListTeamManagersUsers)
				}
				whatsOn := public.Group("/whatsOn")
				{
					whatsOn.GET("", r.afc.WhatsOn.ListAllWhatsOn)
					whatsOn.GET("/latest", r.afc.WhatsOn.GetWhatsOnLatest)
					whatsOn.GET("/past", r.afc.WhatsOn.ListAllWhatsOnEventPast)
					whatsOn.GET("/future", r.afc.WhatsOn.ListAllWhatsOnEventFuture)
					whatsOn.GET("/:id", r.afc.WhatsOn.GetWhatsOnById)
				}
			}
		}
	}
}

func (r *Router) loadBSWDIRoutes() {
	bswdiGroup := r.router.Group("/722e33d4be5c82dfbf45d585679f6c43") // bswdi_group --> MD5
	{
		_ = bswdiGroup
	}
}
