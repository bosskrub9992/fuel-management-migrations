package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         2,
		Up:         up2,
		VerifyUp:   verifyUp2,
		Down:       down2,
		VerifyDown: verifyDown2,
	})
}

func up2(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS cars (
			id BIGSERIAL PRIMARY KEY NOT NULL,
			name VARCHAR(500) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL,
			update_time TIMESTAMP WITH TIME ZONE NOT NULL
		);`,

		`ALTER TABLE fuel_refills 
		ADD COLUMN car_id BIGINT NOT NULL, 
		ADD COLUMN update_by BIGINT NOT NULL,
		ADD COLUMN create_by BIGINT NOT NULL;`,

		`ALTER TABLE fuel_refills DROP COLUMN refill_by;`,

		`ALTER TABLE fuel_usages ADD COLUMN car_id BIGINT NOT NULL;`,

		`ALTER TABLE users ADD COLUMN default_car_id BIGINT NOT NULL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp2(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`SELECT id, name, create_time, update_time FROM cars LIMIT 1;`,
		`SELECT car_id, update_by, create_by FROM fuel_refills LIMIT 1;`,
		`SELECT car_id FROM fuel_usages LIMIT 1;`,
		`SELECT default_car_id FROM users LIMIT 1;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	migrator := tx.Migrator()
	table := "fuel_refills"
	if migrator.HasColumn(table, "refill_by") {
		err := fmt.Errorf("table [%s] is still exists", table)
		slog.Error(err.Error())
		return err
	}
	return nil
}

func down2(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`DROP TABLE IF EXISTS cars;`,

		`ALTER TABLE fuel_refills 
		DROP COLUMN car_id,
		DROP COLUMN update_by,
		DROP COLUMN create_by;`,

		`ALTER TABLE fuel_refills ADD COLUMN refill_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM';`,

		`ALTER TABLE fuel_usages DROP COLUMN car_id;`,
		`ALTER TABLE users DROP COLUMN default_car_id;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown2(ctx context.Context, tx *gorm.DB) error {
	tables := []string{
		"cars",
	}
	migrator := tx.Migrator()
	for index, table := range tables {
		if migrator.HasTable(table) {
			err := fmt.Errorf("table [%s] is still exists", table)
			slog.Error(err.Error(), "index", index)
			return err
		}
	}

	tableToColumns := map[string][]string{
		"fuel_refills": {"car_id", "update_by", "create_by"},
		"fuel_usages":  {"car_id"},
		"users":        {"default_car_id"},
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

	sqlStatements := []string{
		`SELECT refill_by FROM fuel_refills LIMIT 1;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}

	return nil
}
