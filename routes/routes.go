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
		port     string
		version  string
		commit   string
		router   *echo.Echo
		access   *utils.Accesser
		mailer   *utils.Mailer
		afc      *afc.Repos
		admin    *adminPackage.Repo
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
		port:     conf.Port,
		version:  conf.Version,
		commit:   conf.Commit,
		router:   echo.New(),
		access:   conf.Access,
		mailer:   conf.Mailer,
		afc:      conf.AFC,
		admin:    conf.Admin,
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
	r.router.GET("/ping", func(c echo.Context) error {
		resp := map[string]time.Time{"pong": time.Now()}
		return c.JSON(http.StatusOK, resp)
	})

	afcGroup := r.router.Group("/ea231a602d352b2bcc5a2acca6022575") // afc_group --> MD5
    {
        afcV1 := afcGroup.Group("/v1")
        {
            afcV1.GET("/ping", func(c echo.Context) error {
                resp := map[string]time.Time{"pong": time.Now()}
                return c.JSON(http.StatusOK, resp)
            })
            afcV1.PUT("/add", r.afc.Users.AddUser)
            internal := afcV1.Group("/internal")
            {
                if !r.router.Debug {
                    internal.Use(r.access.AFCAuthMiddleware)
                }
                affiliation := internal.Group("/affiliation")
                {
                    _ = affiliation
                }
                document := internal.Group("/document")
                {
                    _ = document
                }
                image := internal.Group("/image")
                {
                    _ = image
                }
                news := internal.Group("/news")
                {
                    _ = news
                }
                player := internal.Group("/player")
                {
                    _ = player
                }
                programme := internal.Group("/programme")
                {
                    _ = programme
                }
                sponsor := internal.Group("/sponsor")
                {
                    _ = sponsor
                }
                team := internal.Group("/team")
                {
                    _ = team
                }
                user := internal.Group("/user")
                {
                    admin := user.Group("/admin")
                    {
                        if !r.router.Debug {
                            admin.Use(r.access.AFCAdminMiddleware)
                        }
                        admin.PUT("", r.afc.Users.AddUser)
                        admin.PATCH("/:email", r.afc.Users.EditUser)
                        admin.DELETE("/:email", r.afc.Users.DeleteUser)
                    }
                    user.GET("/full", r.afc.Users.UserByTokenFull)
                    user.GET("/all", r.afc.Users.ListAllUsers)
                    user.GET("/:email/full", r.afc.Users.UserByEmailFull)
                    user.GET("/:email", r.afc.Users.UserByEmail)
                    user.GET("", r.afc.Users.UserByToken)
                }
                whatsOn := internal.Group("/whatsOn")
                {
                    _ = whatsOn
                }
            }
            public := afcV1.Group("/public")
            {
                public.GET("/affiliations", r.afc.Affiliations.ListAllAffiliations)
                public.GET("/contacts", r.afc.Users.ListAllContactUsers)
                public.GET("/documents", r.afc.Documents.ListAllDocuments)
                public.GET("/images", r.afc.Images.ListAllImages)
                news := public.Group("/news")
                {
                    news.GET("", r.afc.News.ListAllNews)
                    news.GET("/:id", r.afc.News.GetNewsByID)
                    news.GET("/latest", r.afc.News.GetNewsLatest)
                }
                public.GET("/players/:teamID", r.afc.Players.ListAllPlayersByTeamID)
                public.GET("/programmes", r.afc.Programmes.ListAllProgrammes)
                sponsors := public.Group("/sponsors")
                {
                    sponsors.GET("", r.afc.Sponsors.ListALlSponsors)
                    sponsors.GET("/:teamID", r.afc.Sponsors.ListAllSponsorsByTeamID)
                }
                public.GET("/teams", r.afc.Teams.ListAllTeams)
                team := public.Group("/team")
                {
                    team.GET("/:id", r.afc.Teams.GetTeamByID)
                    team.GET("/manager/:id", r.afc.Teams.GetTeamManagerByID)
                }
                whatson := public.Group("/whatsOn")
                {
                    whatson.GET("", r.afc.WhatsOn.ListAllWhatsOn)
                    whatson.GET("/:id", r.afc.WhatsOn.GetWhatsOnByID)
                }
            }
        }
    }

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
		}
	}

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
}
