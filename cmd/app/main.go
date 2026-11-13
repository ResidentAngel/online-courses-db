package main

import (
	"fmt"
	"log"
	"net/http"

	"online-courses-db/internal/config"
	"online-courses-db/internal/handler"
	"online-courses-db/internal/repository"
)

func main() {
	cfg := config.Load()

	store, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Println("Предупреждение БД:", err)
	}

	mux := handler.NewRouter(store)

	fmt.Printf("Сервер запущен на http://localhost%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
