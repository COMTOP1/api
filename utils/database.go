package utils

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

// DatabaseConfig represents a configuration to connect to an SQL database
type DatabaseConfig struct {
	Host     string
	SSLMode  string
	Bucket   string
	Username string
	Password string
}

// NewDB Initialises the connection to the database
func NewDB(config DatabaseConfig) (bucket *gocb.Bucket, err error) {
	cluster, err := gocb.Connect("couchbase://"+config.Host, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: config.Username,
			Password: config.Password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket = cluster.Bucket(config.Bucket)

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		bucket = nil
	}

	return
}
