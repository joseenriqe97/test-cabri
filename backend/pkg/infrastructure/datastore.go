package infrastructure

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joseenriqe97/test-cabri/config"
	"github.com/sirupsen/logrus"
)

func Migrate() {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = formatDir(dir)

	migrations := struct {
		Conn   *string
		Source string
	}{
		Conn:   config.PgConnMigration(),
		Source: "file://" + path.Join(dir, "/pkg/infrastructure/migrations"),
	}

	migration, err := migrate.New(migrations.Source, *migrations.Conn)
	if err != nil {
		logrus.Errorf("error loading migration sourcedir=%s error=%s %v\n", migrations.Source, err.Error(), err)
		panic(err)
	}

	err = migration.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logrus.Errorf("[%s]: Error applying migration. %s. data %v\n", migrations.Source, err.Error(), err)
			panic(err)
		}
		logrus.Infof("[%s]: There are no changes in migrations!\n", migrations.Source)
	}

	migration.Close()

}

func formatDir(dir string) string {
	basedir := strings.Split(dir, "test-cabri")[0]

	return basedir + "/test-cabri/backend"
}
