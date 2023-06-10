package postgresql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Options func(*gorm.Config)

// WithPrepareStatement enable prepared statement for GORM.
//
// Refer to https://gorm.io/docs/session.html#PrepareStmt.
func WithPrepareStatement() Options {
	return func(c *gorm.Config) {
		c.PrepareStmt = true
	}
}

// WithCustomLogger use the given logger implementation as GORM logger.
func WithCustomLogger(l logger.Interface) Options {
	return func(c *gorm.Config) {
		c.Logger = l
	}
}

// WithDisableForeignKeyWhenMigrating disable foreign key constraint when using
// GORM's migration.
func WithDisableForeignKeyWhenMigrating() Options {
	return func(c *gorm.Config) {
		c.DisableForeignKeyConstraintWhenMigrating = true
	}
}

// WithDisableNestedTrx disable default nested transaction.
func WithDisableNestedTrx() Options {
	return func(c *gorm.Config) {
		c.DisableNestedTransaction = true
	}
}

// WithSkipDefaultTrx skip default transaction.
func WithSkipDefaultTrx() Options {
	return func(c *gorm.Config) {
		c.SkipDefaultTransaction = true
	}
}

// WithSingularTableName set all tables name to singular.
func WithSingularTableName() Options {
	return func(c *gorm.Config) {
		c.NamingStrategy = schema.NamingStrategy{
			SingularTable: true,
		}
	}
}
