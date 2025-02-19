package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"

	"fmt"
)

type config struct {
	// Database struct {
	// 	Host     string `env:"DATABASE_HOST,required "`
	// 	Port     int    `env:"DATABASE_PORT,default=5432"`
	// 	User     string `env:"DATABASE_USER,required"`
	// 	Password string `env:"DATABASE_PASSWORD,required"`
	// 	DbName   string `env:"DATABASE_DB_NAME,required"`
	// }
	Server struct {
		HTTPPort int32 `env:"HTTPPORT,default=3001"`
	}

	AuthToken string `env:"AUTH_TOKEN,default=secret"`
}

var c config

func ReadConfig() error {
	ctx := context.Background()
	err := envconfig.Process(ctx, &c)
	return err
}

func HTTPListener() string {
	return fmt.Sprintf(":%d", c.Server.HTTPPort)
}
