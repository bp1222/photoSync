package tinybeans

import (
	"net/http"

	"github.com/bp1222/photoSync/database"
)

type User struct {
	Id       int64    `yaml:"id"`
	FrameIds []string `yaml:"frames"`
}

type Journal struct {
	Id    int64  `yaml:"id"`
	Users []User `yaml:"users"`
}

type Config struct {
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
	Journals []Journal `yaml:"journals"`

	client   *http.Client
	database database.Database
}

type OptionFunc func(*Config)

func WithClient(client *http.Client) OptionFunc {
	return func(c *Config) {
		c.client = client
	}
}

func WithDatabase(db database.Database) OptionFunc {
	return func(c *Config) {
		c.database = db
	}
}
