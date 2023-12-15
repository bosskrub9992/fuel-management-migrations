package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"time"

	"github.com/bosskrub9992/fuel-management-migrations/config"
	"github.com/bosskrub9992/fuel-management-migrations/domains"
	"github.com/bosskrub9992/fuel-management-migrations/migrations/mgsqlite"
	"github.com/bosskrub9992/fuel-management-migrations/slogger"
	"github.com/jinleejun-corp/corelib/databases"
	"gorm.io/gorm"
)

func main() {
	args := os.Args

	var all bool
	if len(args) > 1 && args[1] == "all" {
		all = true
	}

	cfg := config.New()
	ctx := context.Background()
	slog.SetDefault(slogger.New())

	// sort ascending
	sort.SliceStable(mgsqlite.Migrations, func(i, j int) bool {
		return mgsqlite.Migrations[i].ID < mgsqlite.Migrations[j].ID
	})

	idToMigration := make(map[uint]mgsqlite.Migration)
	for _, migration := range mgsqlite.Migrations {
		if _, found := idToMigration[migration.ID]; found {
			slog.Error(fmt.Sprintf("duplicate migration id: [%d]",
				migration.ID,
			))
			return
		}
		idToMigration[migration.ID] = migration
	}

	db, err := databases.NewGormDBSqlite(cfg.Database.FilePath)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	dbMigrator := db.Migrator()
	tableMigration := domains.Migration{}.TableName()
	if !dbMigrator.HasTable(tableMigration) {
		slog.Info(fmt.Sprintf("not found table [%s]", tableMigration))
		if err := dbMigrator.CreateTable(&domains.Migration{}); err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info(fmt.Sprintf("created table [%s]", tableMigration))
	}

	var inDBMigrations []domains.Migration
	if err := db.Model(&domains.Migration{}).Find(&inDBMigrations).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	idToInDBMigration := make(map[uint]domains.Migration)
	for _, migration := range inDBMigrations {
		idToInDBMigration[migration.ID] = migration
	}

	var migratedCount int

	for _, migration := range mgsqlite.Migrations {
		if !all && migratedCount == 1 {
			break
		}
		if _, found := idToInDBMigration[migration.ID]; found {
			continue
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(ctx, tx); err != nil {
				slog.Error(err.Error())
				return err
			}
			if err := migration.VerifyUp(ctx, tx); err != nil {
				slog.Error(err.Error())
				return err
			}
			migrated := domains.Migration{
				ID:        migration.ID,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&migrated).Error; err != nil {
				slog.Error(err.Error())
				return err
			}
			return nil
		})
		if err != nil {
			return
		}
		slog.Info(fmt.Sprintf("succesfully migrated id: [%d] up", migration.ID))
		migratedCount++
	}

	if migratedCount == 0 {
		slog.Info("no up migrations to migrate")
	} else {
		if all {
			slog.Info("successfully migrated up all pending migrations")
		}
	}
}
