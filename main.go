package main

import (
	"log"
	"net/http"
	"task-one/delivery"
	"task-one/pkg/database"
	"task-one/pkg/redis"
	dbRepo "task-one/repository/database"
	redisRepo "task-one/repository/redis"
	"task-one/usecase"

	"github.com/joho/godotenv"
)

func main() {
	// read .env
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load .env")
	}
	//________________

	// Initialize database
	dbRepo.InitDatabase()
	defer database.CloseDB()
	//_________________

	// Initialize Redis
	redisClient := redis.InitRedis()

	// Initialize repositories
	repo := dbRepo.NewRepository()         // postgres
	rdb := redisRepo.NewRedis(redisClient) // redis
	//_________________

	// Initialize usecase
	uc := usecase.NewUsecase(repo, rdb)
	//_________________

	// Initialize router
	router := delivery.NewRouter(uc)

	// Start HTTP server
	log.Println("Starting server on :8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Failed to start server: %+v", err)
	}
	//_________________

	log.Println("Application is running...")
}
