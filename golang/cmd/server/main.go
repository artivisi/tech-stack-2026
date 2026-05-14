package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artivisi/tech-stack-2026/golang/internal/config"
	"github.com/artivisi/tech-stack-2026/golang/internal/repository"
	"github.com/artivisi/tech-stack-2026/golang/internal/web"
)

func main() {
	port := config.RequiredEnv("HTTP_PORT")
	dbURL := config.RequiredEnv("DATABASE_URL")

	db, err := config.ConnectDB(dbURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	tmpl, err := web.NewTemplates()
	if err != nil {
		log.Fatalf("template load: %v", err)
	}

	v := web.NewValidator()
	repo := repository.NewRegistration(db)
	handler := web.NewHandler(repo, tmpl, v)

	mux := http.NewServeMux()
	handler.Register(mux)
	mux.Handle("GET /css/", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           web.RequestLogger(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("registration-golang listening on http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("%s received - shutting down", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
