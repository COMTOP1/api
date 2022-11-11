package routes

import (
	"fmt"
	"github.com/COMTOP1/api/controllers/afc"
	"net/http"
	"time"

	adminPackage "github.com/COMTOP1/api/controllers/admin/v1/admin"
	afcPackages "github.com/COMTOP1/api/controllers/afc"
	"github.com/COMTOP1/api/middleware"
	"github.com/COMTOP1/api/utils"
	"github.com/labstack/echo/v4"
)

type (
	Admin struct {
		Admin *adminPackage.Repo
	}

	AFC struct {
		Users *afcPackages.Repos
	}

	Router struct {
		port     string
		version  string
		commit   string
		router   *echo.Echo
		access   *utils.Accesser
		mailer   *utils.Mailer
		afc      *afc.Repos
		adminURL string
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
		AdminURL   string
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
		adminURL: conf.AdminURL,
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
			internal := afcV1.Group("/internal")
			{
				if !r.router.Debug {
					internal.Use(r.access.AFCAuthMiddleware)
				}
				user := internal.Group("/user")
				{
					user.GET("/full", r.afc.Users.UserByTokenFull)
					user.GET("/:email", r.afc.Users.UserByEmail)
					user.GET("/:email/full", r.afc.Users.UserByEmailFull)
					user.GET("/all", r.afc.Users.ListAllUsers)
					user.GET("", r.afc.Users.UserByToken)
					admin := user.Group("/admin")
					{
						admin.Use(r.access.AFCAdminMiddleware)
					}
				}
			}
			public := afcV1.Group("/public")
			{
				public.GET("/contacts", r.afc.Users.ListAllContactUsers)
			}
		}
	}

	admin := r.router.Group("/" + r.adminURL)
	{
		getJWT := admin.Group("/get")
		{
			getJWT.Use(r.access.AdminInitAuthMiddleware)
			getJWT.GET("/jwt", r.admin.GetJWT)
		}
		sso := admin.Group("/sso")
		{
			sso.Use(r.access.AdminAuthMiddleware)
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
