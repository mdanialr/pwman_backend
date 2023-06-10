package seeder

import (
	"fmt"

	"gorm.io/gorm"
)

// Run inject available seeders to given gorm.DB.
func Run(db *gorm.DB) {
	fmt.Println("Running All Seeders")
	for _, seeder := range seeders {
		fmt.Println("---- Run", seeder.Name, "----")
		seeder.Run(db)
		fmt.Println("---- Done seeding", seeder.Name, "----")
	}
	fmt.Println("Done All Seeders")
}

// seeders hold all seeder that should be run.
var seeders = []struct {
	Name string
	Run  func(*gorm.DB)
}{
	{
		Name: "Registered OTP Entity",
		Run:  registeredOtp,
	},
	{
		Name: "Category Entity",
		Run:  category,
	},
}
