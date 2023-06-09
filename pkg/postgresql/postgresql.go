package postgresql

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGorm init new GORM using postgresql as the db and given viper config to
// get detailed information about the db such as host, port, username etc.
func NewGorm(v *viper.Viper, opts ...Options) (*gorm.DB, error) {
	// DSN format refer to https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s",
		v.GetString("db.host"),
		v.GetInt("db.port"),
		v.GetString("db.name"),
		v.GetString("db.user"),
		v.GetString("db.pass"),
	)
	// prepare gorm config
	conf := gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              false,
		DisableForeignKeyConstraintWhenMigrating: false,
	}
	// apply all provided options
	for _, opt := range opts {
		opt(&conf)
	}
	// open connection to db based on the constructed dsn connection string above
	return gorm.Open(postgres.Open(dsn), &conf)
}
