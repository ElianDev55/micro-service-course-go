package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/ElianDev55/first-api-go/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func DBConnection() (*gorm.DB, error) {

	host := os.Getenv("DB_HOST")
	userDb:= os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	nameDd := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT") 
	debug := os.Getenv("DB_DEBUG") 
	AutoMigrate := os.Getenv("DB_AUTOMIGRATE") 


	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, userDb, password, nameDd, port  )
	db, errDd := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	
	if errDd != nil {
		fmt.Println(errDd)
		return nil, errDd
	}

	if  debug == "true" {
		db = db.Debug()
	}

	if  AutoMigrate == "true" {
		errUser := db.AutoMigrate(&domain.User{})
		if errUser != nil {
			return nil,errUser
		}

		errCour := db.AutoMigrate(&domain.Course{})
		if errCour != nil {
			return nil,errCour
		}
		
		errEnroll := db.AutoMigrate(&domain.Enrollment{})
		if errEnroll != nil {
			return nil,errEnroll
		}

	}

	return db, nil

		
	


}
