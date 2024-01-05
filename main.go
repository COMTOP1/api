package main

import (
	"log"

	"github.com/BurntSushi/toml"

	"github.com/COMTOP1/api/controllers"
	"github.com/COMTOP1/api/controllers/admin/v1/admin"
	"github.com/COMTOP1/api/controllers/afc"
	"github.com/COMTOP1/api/controllers/bswdi"
	"github.com/COMTOP1/api/controllers/sso"
	"github.com/COMTOP1/api/handler"
	"github.com/COMTOP1/api/routes"
	"github.com/COMTOP1/api/structs"
	"github.com/COMTOP1/api/utils"
)

func main() {
	config := &structs.Config{}
	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Server.Debug {
		log.SetFlags(log.Llongfile)
	}

	//var mailer *utils.Mailer
	//if config.Mail.Host != "" {
	//	if config.Mail.Enabled {
	//		mailConfig := utils.MailConfig{
	//			Host:     config.Mail.Host,
	//			Port:     config.Mail.Port,
	//			Username: config.Mail.User,
	//			Password: config.Mail.Password,
	//		}
	//
	//		mailer, err = utils.NewMailer(mailConfig)
	//		if err != nil {
	//			log.Printf("failed to connect to mail server: %+v", err)
	//			config.Mail.Enabled = false
	//		} else {
	//			log.Printf("Connected to mail server: %s\n", config.Mail.Host)
	//
	//			mailer.Defaults = utils.Defaults{
	//				DefaultTo:   "root@bswdi.co.uk",
	//				DefaultFrom: "BSWDI API <api@bswdi.co.uk>",
	//			}
	//		}
	//	}
	//} else {
	//	config.Mail.Enabled = false
	//}

	dbConfig := utils.DatabaseConfig{
		Host:     config.Database.Host,
		SSLMode:  config.Database.SSLMode,
		Bucket:   config.Database.Bucket,
		Username: config.Database.User,
		Password: config.Database.Password,
	}
	bucket, err := utils.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("failed to start DB: %+v", err)
	}
	log.Printf("Connected to DB: %s@%s", dbConfig.Username, dbConfig.Host)

	access := utils.NewAccesser(utils.AccesserConfig{
		AccessCookieName: config.Server.Access.AccessCookieName, // jwt_token --> base64
		SigningKey:       []byte(config.Server.Access.SigningToken),
		DomainName:       config.Server.DomainName,
		Admin: struct {
			AdminAccessCookieName string
			Key0                  string
			Key1                  string
			Key2                  string
			Key3                  string
			TOTP                  string
		}{
			AdminAccessCookieName: config.Server.Admin.AdminAccessCookieName,
			Key0:                  config.Server.Admin.Key0,
			Key1:                  config.Server.Admin.Key1,
			Key2:                  config.Server.Admin.Key2,
			Key3:                  config.Server.Admin.Key3,
			TOTP:                  config.Server.Admin.TOTP,
		},
	})

	log.Printf("API Version %s", config.Server.Version)

	afcScope := bucket.Scope("afc")

	adminScope := bucket.Scope("admin")

	bswdiScope := bucket.Scope("bswdi")

	ssoScope := bucket.Scope("sso")

	session, err := handler.NewSession(config.Server.DomainName)
	if err != nil {
		log.Fatalf("The session couldn't be initialised!\n\n%s\n\nExiting!", err)
	}

	controller := controllers.GetController(access, session)

	router := routes.New(&routes.NewRouter{
		Port:       config.Server.Port,
		Version:    config.Server.Version,
		Commit:     config.Server.Commit,
		DomainName: config.Server.DomainName,
		Debug:      config.Server.Debug,
		Access:     access,
		Mailer: utils.NewMailer(utils.Config{
			Host:     config.Mail.Host,
			Port:     config.Mail.Port,
			Username: config.Mail.User,
			Password: config.Mail.Password,
		}),
		AFC:   afc.NewRepos(afcScope, controller, config.Server.ServiceURL.AFC),
		Admin: admin.NewRepo(adminScope, controller, config.Server.DomainName, config.Server.Admin.URL, []byte(config.Server.Access.SigningToken)),
		BSWDI: bswdi.NewRepo(bswdiScope, controller, config.Server.ServiceURL.BSWDI),
		SSO:   sso.NewRepo(ssoScope, controller, config.Server.ServiceURL.SSO),
	})

	err = router.Start()
	if err != nil {
		log.Fatalf("The web server couldn't be started!\n\n%s\n\nExiting!", err)
	}
}
