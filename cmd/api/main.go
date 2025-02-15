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
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: error %s \n", err)
	}

	handler := api.NewHandler(db)

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
