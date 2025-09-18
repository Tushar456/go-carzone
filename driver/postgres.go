package driver

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	constStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))

	//db, err := sql.Open("postgres", constStr)
	db, err := gorm.Open(postgres.Open(constStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error %s when opening DB\n", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error %s when getting generic DB\n", err)
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Error %s pinging DB\n", err)
		return nil, err
	}

	fmt.Println("Connected to DB successfully")

	return db, nil
}
