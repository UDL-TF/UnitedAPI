package context

import (
	"context"

	"gorm.io/gorm"
)

// AppContext holds all application dependencies and services
// This is the central dependency injection container
type AppContext struct {
	// Database connection (example - add when needed)
	DB *gorm.DB
	// Context for cancellation/timeout
	Ctx context.Context
}

// NewAppContext creates a new application context
func NewAppContext(ctx context.Context) *AppContext {
	return &AppContext{
		Ctx: ctx,
	}
}

func (ac *AppContext) SetDB(db *gorm.DB) *AppContext {
	ac.DB = db
	return ac
}

func (ac *AppContext) GetDB() *gorm.DB {
	return ac.DB
}
