package migration

import (
	"fmt"
	"log"
	"os"

	conf "github.com/mdanialr/pwman_backend/pkg/config"
	gl "github.com/mdanialr/pwman_backend/pkg/gorm"
	"github.com/mdanialr/pwman_backend/pkg/migration/seeder"
	"github.com/mdanialr/pwman_backend/pkg/postgresql"

	"gorm.io/gorm"
)

// Run do run migration (creating all tables) and optionally run seeder if
// given param is true.
func Run(isSeeder, isDrop bool) {
	db := initGorm()
	// get the sql db
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("failed to get the DB instance from gorm:", err)
		return
	}
	defer sqlDB.Close()

	// drop tables
	if isDrop {
		fmt.Println("Dropping All Tables")
		db.Migrator().DropTable()
		fmt.Println("Done Dropping All Tables")
	}

	// create tables
	fmt.Println("Creating All Tables")
	db.Migrator().AutoMigrate()
	fmt.Println("Done Creating All Tables")

	// seed the tables with fake data from seeders
	if isSeeder {
		seeder.Run(db)
	}
}

func initGorm() *gorm.DB {
	// init viper config
	v, err := conf.InitConfigYml()
	if err != nil {
		log.Fatalln("failed to init config:", err)
	}
	// setup the logger for GORM
	gormLog := gl.New(os.Stdout,
		gl.WithLogLevel(2),
		gl.WithThreshold(150),
		gl.WithIgnoreRecordNotFound(),
	)
	// init gorm using postgresql as the DB
	db, err := postgresql.NewGorm(v,
		postgresql.WithCustomLogger(gormLog),
		postgresql.WithPrepareStatement(),
		postgresql.WithDisableForeignKeyWhenMigrating(),
	)
	if err != nil {
		log.Fatalln("failed to init gorm with mysql as the DB:", err)
		return nil
	}

	return db
}
