package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rbcervilla/redisstore/v8"
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

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		panic("failed to create redis store")
	}

	store.KeyPrefix("session_")
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	e := api.NewRouter(db, store)
	e.Logger.Debug(e.Start(":8080"))
}
