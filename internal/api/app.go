package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/handlers"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/services"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/config"
	"github.com/vseriousv/go_diffie_hellman_chat_server/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AppStruct struct {
	Router chi.Router
}

func (a *AppStruct) RunApp(c *config.Config, log *zap.Logger) {
	// DATABASE
	db, err := connectDB(c.DataBaseUrl)
	utils.HandleError(err, "DB is not connection")
	err = db.Ping(context.Background())
	utils.HandleError(err, "DB is not connection")
	defer db.Close()

	router := chi.NewRouter()
	router.Use(requestLogger(log))

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// services with db
	messageService := services.NewMessageService(db)

	// handlers
	messageHandler := handlers.NewMessageHandler(messageService)

	router.Post("/v1/messages", messageHandler.CreateMessage)
	router.Get("/v1/messages/{publicKey}", messageHandler.GetMessagesByPublicKey)

	log.Info(fmt.Sprintf("API server is listening on http://%s", c.RunApp))
	if err := http.ListenAndServe(c.RunApp, router); err != nil {
		log.Error(err.Error())
	}
}

func connectDB(connectionString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	poolConfig.MinConns = 2
	poolConfig.MaxConns = 8

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func requestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info("",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Duration("time", time.Since(start)),
			)
		})
	}
}
