package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         5,
		Up:         up5,
		VerifyUp:   verifyUp5,
		Down:       down5,
		VerifyDown: verifyDown5,
	})
}

func up5(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_usages DROP COLUMN is_paid;`,
		`ALTER TABLE fuel_usage_users ADD COLUMN is_paid BOOL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp5(ctx context.Context, tx *gorm.DB) error {
	stmt := `SELECT is_paid FROM fuel_usage_users LIMIT 1;`
	if err := tx.WithContext(ctx).Exec(stmt).Error; err != nil {
		slog.Error(err.Error())
		return err
	}
	migrator := tx.Migrator()
	if migrator.HasColumn("fuel_usages", "is_paid") {
		err := fmt.Errorf("table [%s] still has field [%s]", "fuel_usages", "is_paid")
		slog.Error(err.Error())
		return err
	}
	return nil
}

func down5(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_usages ADD COLUMN is_paid BOOL;`,
		`ALTER TABLE fuel_usage_users DROP COLUMN is_paid;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown5(ctx context.Context, tx *gorm.DB) error {
	stmt := `SELECT is_paid FROM fuel_usages LIMIT 1;`
	if err := tx.WithContext(ctx).Exec(stmt).Error; err != nil {
		slog.Error(err.Error())
		return err
	}
	migrator := tx.Migrator()
	if migrator.HasColumn("fuel_usage_users", "is_paid") {
		err := fmt.Errorf("table [%s] still has field [%s]", "fuel_usage_users", "is_paid")
		slog.Error(err.Error())
		return err
	}
	return nil
}
