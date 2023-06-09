package gorm_logger

import (
	"io"
	"log"
	"time"

	"gorm.io/gorm/logger"
)

// New use given writer as the target where the log for GORM is written to.
// Also apply the provided various available options.
func New(wr io.Writer, opts ...Options) logger.Interface {
	cnf := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: false,
		Colorful:                  false,
	}

	for _, opt := range opts {
		opt(&cnf)
	}

	newGormLogger := logger.New(
		log.New(wr, "\r\n", log.LstdFlags),
		cnf,
	)
	return newGormLogger
}
