package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type config struct {
	Database struct {
		Host     string `env:"DATABASE_HOST,required"`
		Port     int    `env:"DATABASE_PORT,default=5432"`
		User     string `env:"DATABASE_USER,required"`
		Password string `env:"DATABASE_PASSWORD,required"`
		DbName   string `env:"DATABASE_DB_NAME,required"`
	}
	Server struct {
		HTTPPort int32 `env:"HTTPPORT,default=3001"`
	}
	Migrate   bool   `env:"MIGRATE,default=true"`
	AuthToken string `env:"AUTH_TOKEN,default=secret"`
	UrlPLD    string `env:"URL_PLD,default=http://98.81.235.22"`
}

var c config

func ReadConfig() error {
	ctx := context.Background()
	err := envconfig.Process(ctx, &c)
	return err
}

func PgConn() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DbName)
}
func PgConnMigration() *string {
	if c.Migrate {
		pgconn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DbName)
		return &pgconn
	}
	return nil
}

func EnableMigrations() bool {
	return c.Migrate
}

func GetSecretJWT() string {
	return c.AuthToken
}

func HTTPListener() string {
	return fmt.Sprintf(":%d", c.Server.HTTPPort)
}
func GetUrlPLD() string {
	return c.UrlPLD
}
