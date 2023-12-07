package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         1,
		Up:         up1,
		VerifyUp:   verifyUp1,
		Down:       down1,
		VerifyDown: verifyDown1,
	})
}

func up1(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS fuel_usages (
			id BIGSERIAL PRIMARY KEY,
			fuel_use_time TIMESTAMP WITH TIME ZONE NOT NULL,
			fuel_price DECIMAL(10,3) NOT NULL,
			kilometer_before_use INT NOT NULL,
			kilometer_after_use INT NOT NULL,
			description VARCHAR(500),
			total_money DECIMAL(10,3) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL,
			update_time TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY NOT NULL,
			nickname VARCHAR(500) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL,
			update_time TIMESTAMP WITH TIME ZONE NOT null
		);`,
		`CREATE TABLE IF NOT EXISTS fuel_usage_users (
			id BIGSERIAL PRIMARY KEY NOT NULL,
			fuel_usage_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS fuel_refills (
			id BIGSERIAL PRIMARY KEY NOT NULL,
			refill_date TIMESTAMP WITH TIME ZONE NOT NULL,
			total_money DECIMAL(10,3) NOT NULL,
			kilometer_before_refill INT NOT NULL,
			kilometer_after_refill INT NOT NULL,
			fuel_price_calculated DECIMAL(10,3) NOT NULL,
			refill_by VARCHAR(50) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL,
			update_time TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp1(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`SELECT * FROM fuel_usages LIMIT 1;`,
		`SELECT * FROM users LIMIT 1;`,
		`SELECT * FROM fuel_usage_users LIMIT 1;`,
		`SELECT * FROM fuel_refills LIMIT 1;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func down1(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`DROP TABLE IF EXISTS fuel_usages;`,
		`DROP TABLE IF EXISTS users;`,
		`DROP TABLE IF EXISTS fuel_usage_users;`,
		`DROP TABLE IF EXISTS fuel_refills;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown1(ctx context.Context, tx *gorm.DB) error {
	tables := []string{
		"fuel_usages",
		"users",
		"fuel_usage_users",
		"fuel_refills",
	}
	migrator := tx.Migrator()
	for index, table := range tables {
		if migrator.HasTable(table) {
			err := fmt.Errorf("table [%s] is still exists", table)
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}
