package migration

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
)

func ShouldMigrateDB() bool {
	flag.Parse()

	return flag.Arg(0) == "migrate"
}

func is_up() bool {
	if flag.Arg(1) == "up" {
		return true
	} else if flag.Arg(1) == "down" {
		return false
	} else {
		panic("Invalid migration direction - must be 'up' or 'down'. ex: `go run cmd/app/main.go migrate up`")
	}
}

func MigrateDB(config *config.Config, ctx context.Context, repo repository.DBTX) {
	_, filename, _, _ := runtime.Caller(0)
	migrations_dir := filepath.Join(filepath.Dir(filename), "..", "..", "database", "migrations")

	filenames, err := os.ReadDir(migrations_dir)
	if err != nil {
		panic(err)
	}

	for _, filename := range filenames {
		content, err := os.ReadFile(filepath.Join(migrations_dir, filename.Name()))
		if err != nil {
			panic(err)
		}

		cmd := ""
		if is_up() && strings.HasSuffix(filename.Name(), "up.sql") {
			cmd += string(content)
		} else if !is_up() && strings.HasSuffix(filename.Name(), "down.sql") && !strings.Contains(filename.Name(), "user") {
			cmd += string(content)
		}

		for _, command := range strings.Split(cmd, ";") {
			tag, _ := repo.Query(ctx, command+";")
			err = tag.Err()
			tag.Close()

			if err != nil {
				panic(command + "\n" + err.Error())
			} else {
				fmt.Println(command)
			}
		}
	}
}
