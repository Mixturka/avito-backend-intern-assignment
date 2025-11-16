package e2e_test

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
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	apiServer  *handlers.ApiV1
	testRouter http.Handler
)

func TestMain(m *testing.M) {
	cfg := config.Config{
		DB: config.DbConfig{
			Name:     "testdb",
			User:     "postgres",
			Password: "postgres",
			Host:     "localhost",
			Port:     "5432",
			SSLMode:  "disable",
		},
		ServerPort: "8081",
	}

	pool, err := pgxpool.New(context.Background(), cfg.DB.URL())
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}
	dbAdapter := &pgxadapter.PoolAdapter{Pool: pool}

	uRepo := userRepo.NewPostgresRepository(dbAdapter)
	tRepo := teamRepo.NewPostgresRepository(dbAdapter)
	prRepo := prRepo.NewPostgresRepository(dbAdapter)

	uService := user.NewService(uRepo)
	tService := team.NewService(tRepo, uRepo, dbAdapter)
	prService := pullrequest.NewService(prRepo, uRepo, dbAdapter)

	prh := prHandler.NewHandler(prService)
	uh := uHandler.NewHandler(uService)
	th := tHandler.NewHandler(tService)

	apiServer = handlers.NewApiV1(th, uh, prh)

	r := chi.NewRouter()
	apiHandler := api.NewStrictHandler(apiServer, nil)
	r.Mount("/", api.Handler(apiHandler))
	testRouter = r

	code := m.Run()
	os.Exit(code)
}
