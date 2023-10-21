package migrations

import (
	"context"

	"gorm.io/gorm"
)

type Migration struct {
	ID         uint
	Up         func(ctx context.Context, tx *gorm.DB) error
	VerifyUp   func(ctx context.Context, tx *gorm.DB) error
	Down       func(ctx context.Context, tx *gorm.DB) error
	VerifyDown func(ctx context.Context, tx *gorm.DB) error
}

var Migrations = []Migration{}
