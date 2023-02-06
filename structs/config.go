package structs

// Config is a structure containing global website configuration.
//
// See the comments for Server and PageContext for more details.
type (
	Config struct {
		Server   Server   `toml:"server"`
		Database Database `toml:"database"`
		Mail     Mail     `toml:"mail"`
	}

	Server struct {
		Debug      bool   `toml:"api_debug"`
		Port       string `toml:"api_port"`
		DomainName string `toml:"api_domain_name"`
		Version    string `toml:"api_version"`
		Commit     string `toml:"api_commit"`
		Access     Access `toml:"access"`
		Admin      Admin  `toml:"admin"`
	}

	Access struct {
		AccessCookieName string `toml:"api_access_cookie_name"`
		SigningToken     string `toml:"api_signing_token"`
	}

	Admin struct {
		AdminAccessCookieName string `toml:"api_admin_access_cookie_name"`
		Key0                  string `toml:"api_admin_key_0"`
		Key1                  string `toml:"api_admin_key_1"`
		Key2                  string `toml:"api_admin_key_2"`
		Key3                  string `toml:"api_admin_key_3"`
		TOTP                  string `toml:"api_admin_totp"`
		URL                   string `toml:"api_admin_url"`
	}

	Database struct {
		Host     string `toml:"api_db_host"`
		Bucket   string `toml:"api_db_bucket"`
		User     string `toml:"api_db_user"`
		Password string `toml:"api_db_pass"`
		SSLMode  string `toml:"api_db_sslmode"`
	}

	Mail struct {
		Enabled  bool   `toml:"api_mail_enabled"`
		Host     string `toml:"api_mail_host"`
		User     string `toml:"api_mail_user"`
		Password string `toml:"api_mail_pass"`
		Port     int    `toml:"api_mail_port"`
	}
)
