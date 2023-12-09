package migrations

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         6,
		Up:         up6,
		VerifyUp:   verifyUp6,
		Down:       down6,
		VerifyDown: verifyDown6,
	})
}

func up6(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills RENAME COLUMN refill_date TO refill_time;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp6(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	if err := columnShouldExist(migrator, "fuel_refills", "refill_time"); err != nil {
		slog.Error(err.Error())
		return err
	}
	if err := columnShouldNotExist(migrator, "fuel_refills", "refill_date"); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func down6(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills RENAME COLUMN refill_time TO refill_date;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown6(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	if err := columnShouldExist(migrator, "fuel_refills", "refill_date"); err != nil {
		slog.Error(err.Error())
		return err
	}
	if err := columnShouldNotExist(migrator, "fuel_refills", "refill_time"); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
