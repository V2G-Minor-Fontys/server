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

func MigrateDB(config *config.Config, ctx context.Context, repo repository.DBTX) {
	_, filename, _, _ := runtime.Caller(0)
	migrations_dir := filepath.Join(filepath.Dir(filename), "..", "..", "database", "migrations")

	filenames, err := os.ReadDir(migrations_dir)
	if err != nil {
		panic(err)
	}

	sqlUpStatement := ""
	sqlDownStatement := ""
	for _, filename := range filenames {
		content, err := os.ReadFile(filepath.Join(migrations_dir, filename.Name()))
		if err != nil {
			panic(err)
		}

		if strings.HasSuffix(filename.Name(), "up.sql") {
			sqlUpStatement += string(content)
		} else if strings.HasSuffix(filename.Name(), "down.sql") {
			sqlDownStatement += string(content)
		}
	}

	tag, _ := repo.Query(ctx, sqlUpStatement)
	fmt.Println(tag.CommandTag().Delete())
	fmt.Println(tag.CommandTag().RowsAffected())
	tag.Close()

	// repo.Exec(ctx, sqlUpStatement)

	// psql "host=localhost port=5432 user=testUser password=TestPass1234! dbname=testDb sslmode=disable" -c (cat database/migrations/000001_create_identities_table.up.sql | tr '\n' ' ')
}
