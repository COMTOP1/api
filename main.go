package main

import (
	"encoding/base64"
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
	hash := []uint8{161, 64, 116, 252, 86, 93, 179, 86, 33, 212, 219, 137, 62, 118, 72, 150, 161, 16, 41, 158, 9, 170, 231, 147, 39, 215, 61, 25, 137, 38, 105, 49, 145, 34, 1, 168, 110, 16, 227, 182, 235, 45, 220, 148, 217, 173, 182, 122, 241, 0, 95, 53, 84, 162, 118, 87, 120, 142, 231, 98, 29, 49, 63, 32}
	salt := []uint8{174, 254, 21, 182, 67, 162, 129, 237, 136, 35, 243, 191, 209, 163, 18, 40, 65, 199, 26, 74, 42, 251, 137, 195, 142, 28, 193, 207, 151, 108, 15, 45, 60, 239, 198, 128, 248, 149, 250, 250, 123, 185, 83, 29, 87, 133, 160, 83, 173, 180, 131, 14, 24, 115, 21, 76, 172, 64, 35, 213, 142, 134, 208, 106}

	hashEncode := base64.StdEncoding.EncodeToString(hash)
	fmt.Println(hashEncode)

	saltEncode := base64.StdEncoding.EncodeToString(salt)
	fmt.Println(saltEncode)

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
