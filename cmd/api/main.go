package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"authorization-vendor/internal/httpserver"
	"authorization-vendor/internal/controller"
	"authorization-vendor/internal/repository/memory"
	"authorization-vendor/internal/service"
)

func main() {
	repo := memory.NewJSONRepository("data/vendedores.json", "data/historial.json")
	svc := service.NewAuthorizationVendorService(repo)
	ctrl := controller.NewAuthorizationVendorController(svc)

	r := httpserver.NewRouter(ctrl)

	addr := getEnv("HTTP_ADDR", ":8081")
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("HTTP server listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
