package initializers

import (
	"ApiGin/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection() {
	var error error
	var DSN = os.Getenv("DB_URL")
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Println("DSN is: ", DSN)
		log.Fatal(error)
	} else {
		log.Println("DB Connected")
	}
	runMigrations(DB)
	log.Println("DSN is: ", DSN)
	log.Println("Migrations Done")
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(

		&models.Employee{},
		&models.EmployeeSkills{},
		&models.Skill{},
		&models.Level{},
		&models.User{},
	)
}
