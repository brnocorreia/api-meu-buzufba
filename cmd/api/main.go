package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brnocorreia/api-meu-buzufba/internal/api"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/database/mongodb"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/database/redis"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/mail"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	ctx := context.Background()

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: error %s \n", err)
	}
	defer db.Client().Disconnect(ctx)

	redis, err := redis.NewRedisConnection(ctx)
	if err != nil {
		log.Fatalf("Error connecting to Redis: error %s \n", err)
	}
	defer redis.Close()

	mail.InitMailer(os.Getenv("RESEND_API_KEY"))

	handler := api.NewHandler(db, redis)

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Error starting server: error %s \n", err)
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")
}
