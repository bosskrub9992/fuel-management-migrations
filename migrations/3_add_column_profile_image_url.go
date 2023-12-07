package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         3,
		Up:         up3,
		VerifyUp:   verifyUp3,
		Down:       down3,
		VerifyDown: verifyDown3,
	})
}

func up3(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE users ADD COLUMN profile_image_url VARCHAR(1000);`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp3(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`SELECT profile_image_url FROM users LIMIT 1;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func down3(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE users DROP COLUMN profile_image_url;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown3(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	tableToColumns := map[string][]string{
		"users": {"profile_image_url"},
	}
	for table, columns := range tableToColumns {
		for _, column := range columns {
			if migrator.HasColumn(table, column) {
				err := fmt.Errorf("table [%s] still has field [%s]", table, column)
				slog.Error(err.Error())
				return err
			}
		}
	}
	return nil
}
