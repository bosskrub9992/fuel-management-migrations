package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/bosskrub9992/fuel-management-migrations/config"
	"github.com/bosskrub9992/fuel-management-migrations/domains"
	"github.com/bosskrub9992/fuel-management-migrations/migrations"
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

	// sort descending
	sort.SliceStable(migrations.Migrations, func(i, j int) bool {
		return migrations.Migrations[i].ID > migrations.Migrations[j].ID
	})

	idToMigration := make(map[uint]migrations.Migration)
	for _, migration := range migrations.Migrations {
		if _, found := idToMigration[migration.ID]; found {
			slog.Error(fmt.Sprintf("duplicate migration id: [%d]",
				migration.ID,
			))
			return
		}
		idToMigration[migration.ID] = migration
	}

	sqlDB, err := databases.NewPostgres(&cfg.Database)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer sqlDB.Close()
	db, err := databases.NewGormDBPostgres(sqlDB)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	dbMigrator := db.Migrator()
	tableMigration := domains.Migration{}.TableName()
	if !dbMigrator.HasTable(tableMigration) {
		slog.Info(fmt.Sprintf("not found table [%s]",
			tableMigration,
		))
		if err := dbMigrator.CreateTable(&domains.Migration{}); err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info(fmt.Sprintf("created table [%s]",
			tableMigration,
		))
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
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, migration := range migrations.Migrations {
			if !all && migratedCount == 1 {
				break
			}
			if _, found := idToInDBMigration[migration.ID]; !found {
				continue
			}
			if err := migration.Down(ctx, tx); err != nil {
				slog.Error(err.Error())
				return err
			}
			if err := migration.VerifyDown(ctx, tx); err != nil {
				slog.Error(err.Error())
				return err
			}
			if err := tx.Where("id = ?", migration.ID).Delete(&domains.Migration{}).Error; err != nil {
				slog.Error(err.Error())
				return err
			}
			slog.Info(fmt.Sprintf("succesfully migrated id: [%d] down",
				migration.ID,
			))
			migratedCount++
		}
		return nil
	})
	if err != nil {
		return
	}

	if migratedCount == 0 {
		slog.Info("no down migrations to migrate")
	} else {
		if all {
			slog.Info("successfully migrated down all migrations")
		}
	}
}
