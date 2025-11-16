package main

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/api/handlers"
	prHandler "avito-backend-intern-assignment/internal/app/api/handlers/pullrequest"
	tHandler "avito-backend-intern-assignment/internal/app/api/handlers/team"
	uHandler "avito-backend-intern-assignment/internal/app/api/handlers/user"
	"avito-backend-intern-assignment/internal/app/application/service/pullrequest"
	"avito-backend-intern-assignment/internal/app/application/service/team"
	"avito-backend-intern-assignment/internal/app/application/service/user"
	prRepo "avito-backend-intern-assignment/internal/app/infrastructure/repository/pullrequest"
	teamRepo "avito-backend-intern-assignment/internal/app/infrastructure/repository/team"
	userRepo "avito-backend-intern-assignment/internal/app/infrastructure/repository/user"
	"avito-backend-intern-assignment/internal/pkg/config"
	"avito-backend-intern-assignment/pkg/db/pgxadapter"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DB.URL())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	dbAdapter := &pgxadapter.PoolAdapter{Pool: pool}

	userRepo := userRepo.NewPostgresRepository(dbAdapter)
	teamRepo := teamRepo.NewPostgresRepository(dbAdapter)
	prRepo := prRepo.NewPostgresRepository(dbAdapter)

	userService := user.NewService(userRepo)
	teamService := team.NewService(teamRepo, userRepo, dbAdapter)
	prService := pullrequest.NewService(prRepo, userRepo, dbAdapter)

	prh := prHandler.NewHandler(prService)
	uh := uHandler.NewHandler(userService)
	th := tHandler.NewHandler(teamService)

	h := handlers.NewApiV1(th, uh, prh)
	r := chi.NewRouter()

	apiHandler := api.NewStrictHandler(h, nil)

	r.Mount("/", api.Handler(apiHandler))

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}
