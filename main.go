package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/COMTOP1/api/controllers/admin/v1/admin"
	"github.com/COMTOP1/api/controllers/afc"
	"github.com/COMTOP1/api/routes"
	"github.com/COMTOP1/api/structs"
	"github.com/COMTOP1/api/utils"
	"html/template"
	"log"
	"os"
	"os/signal"
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

	var mailer *utils.Mailer
	if config.Mail.Host != "" {
		mailConfig := utils.MailConfig{
			Host:     config.Mail.Host,
			Port:     config.Mail.Port,
			Username: config.Mail.User,
			Password: config.Mail.Password,
		}

		mailer, err = utils.NewMailer(mailConfig)
		if err != nil {
			log.Printf("failed to connect to mail server: %+v", err)
		}

		mailer.Defaults = utils.Defaults{
			DefaultTo:   "root@bswdi.co.uk",
			DefaultFrom: "BSWDI API <api@bswdi.co.uk>",
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if mailer != nil {
				exitingTemplate := template.New("Exiting Template")
				exitingTemplate = template.Must(exitingTemplate.Parse("<body>BSWDI API has been stopped!<br><br>{{if .Debug}}Exit signal: {{.Sig}}<br><br>{{end}}Version: {{.Version}}<br>Commit: {{.Commit}}</body>"))

				starting := utils.Mail{
					Subject:     "BSWDI API has been stopped!",
					UseDefaults: true,
					Tpl:         exitingTemplate,
					TplData: struct {
						Debug           bool
						Sig             os.Signal
						Version, Commit string
					}{
						Debug:   config.Server.Debug,
						Sig:     sig,
						Version: config.Server.Version,
						Commit:  config.Server.Commit,
					},
				}

				err = mailer.SendMail(starting)
				if err != nil {
					fmt.Println(err)
				}
				err = mailer.Close()
				if err != nil {
					fmt.Println(err)
				}
			}
			os.Exit(0)
		}
	}()

	dbConfig := utils.DatabaseConfig{
		Host:     config.Database.Host,
		SSLMode:  config.Database.SSLMode,
		Bucket:   config.Database.Bucket,
		Username: config.Database.User,
		Password: config.Database.Password,
	}
	bucket, err := utils.NewDB(dbConfig)
	if err != nil {
		if mailer != nil {
			err1 := mailer.SendErrorFatalMail(utils.Mail{
				Error:       fmt.Errorf("failed to start DB: %+v", err),
				UseDefaults: true,
			})
			if err1 != nil {
				fmt.Println(err1)
			}
		}
		log.Fatalf("failed to start DB: %+v", err)
	}
	log.Printf("Connected to DB: %s@%s", dbConfig.Username, dbConfig.Host)

	access := utils.NewAccesser(utils.Config{
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

	if mailer != nil {

		startingTemplate := template.New("Starting Template")
		startingTemplate = template.Must(startingTemplate.Parse("<body>BSWDI API starting{{if .Debug}} in debug mode!<br><b>Do not run in production! Authentication is disabled!</b>{{else}}!{{end}}<br><br>Version: {{.Version}}<br>Commit: {{.Commit}}<br><br>If you don't get another email then this has started correctly.</body>"))

		subject := "BSWDI API is starting"

		if config.Server.Debug {
			subject += " in debug mode"
			log.Println("Debug Mode - Disabled auth - do not run in production!")
		}

		subject += "!"

		starting := utils.Mail{
			Subject:     subject,
			UseDefaults: true,
			Tpl:         startingTemplate,
			TplData: struct {
				Debug           bool
				Version, Commit string
			}{
				Debug:   config.Server.Debug,
				Version: config.Server.Version,
				Commit:  config.Server.Commit,
			},
		}

		err = mailer.SendMail(starting)
		if err != nil {
			fmt.Println(err)
		}
	}

	afcScope := bucket.Scope("afc")

	adminScope := bucket.Scope("admin")

	err = routes.New(&routes.NewRouter{
		Port:       config.Server.Port,
		Version:    config.Server.Version,
		Commit:     config.Server.Commit,
		DomainName: config.Server.DomainName,
		Debug:      config.Server.Debug,
		Access:     access,
		Mailer:     mailer,
		AFC:        afc.NewRepos(afcScope, access),
		Admin:      admin.NewRepo(adminScope, access, config.Server.DomainName, config.Server.Admin.URL, config.Server.Access.SigningToken),
	}).Start()
	if err != nil {
		if mailer != nil {
			err1 := mailer.SendErrorFatalMail(utils.Mail{
				Error:       fmt.Errorf("the web server couldn't be started: %s... exiting", err),
				UseDefaults: true,
			})
			if err1 != nil {
				fmt.Println(err1)
			}
		}
		log.Fatalf("The web server couldn't be started!\n\n%s\n\nExiting!", err)
	}
}
