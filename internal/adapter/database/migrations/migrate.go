package migrations

import (
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.sql
var sqlFiles embed.FS

func RunMigrations(db *gorm.DB) error {
	entries, err := sqlFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting sql.DB: %w", err)
	}

	for _, name := range names {
		content, err := sqlFiles.ReadFile(name)
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}
		if _, err := sqlDB.Exec(string(content)); err != nil {
			return fmt.Errorf("executing %s: %w", name, err)
		}
		log.Printf("migration applied: %s", name)
	}

	return nil
}
