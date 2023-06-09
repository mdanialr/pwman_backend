package gorm_logger

import (
	"time"

	"gorm.io/gorm/logger"
)

type Options func(*logger.Config)

// WithThreshold set the duration of GORM's operation that's regarded as slow
// operation in milliseconds.
func WithThreshold(dur int) Options {
	return func(c *logger.Config) {
		c.SlowThreshold = time.Duration(dur) * time.Millisecond
	}
}

// WithIgnoreRecordNotFound shutdown the error record not found, because by
// default GORM will throw error on record not found.
func WithIgnoreRecordNotFound() Options {
	return func(c *logger.Config) {
		c.IgnoreRecordNotFoundError = true
	}
}

// WithLogLevel set the GORM's logger log level. Should be set as follows
//
//	1: Silent
//	2: Error (only log error query)
//	3: Warn (include warning query)
//	4: Info (most verbose)
func WithLogLevel(lvl int) Options {
	return func(c *logger.Config) {
		switch lvl {
		case 1:
			c.LogLevel = logger.Silent
		case 2:
			c.LogLevel = logger.Error
		case 3:
			c.LogLevel = logger.Warn
		case 4:
			c.LogLevel = logger.Info
		}
	}
}

// WithColor set the GORM log to be colorful.
func WithColor() Options {
	return func(c *logger.Config) {
		c.Colorful = true
	}
}
