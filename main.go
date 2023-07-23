package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/regmarmcem/echo-session-demo/api"
	"github.com/regmarmcem/echo-session-demo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}

	db.AutoMigrate(&model.User{})
	e := api.NewRouter(db)
	e.Logger.Panic(e.Start(":8080"))
}
