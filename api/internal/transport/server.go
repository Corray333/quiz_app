package transport

import (
	"log/slog"
	"net/http"

	"github.com/Corray333/quiz/internal/types"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Question interface {
	GetType() string
}

type Service interface {
	CreateQuiz(quiz *types.Quiz) (int64, error)
	CreateQuestion(question *types.Question) (int64, error)
	GetQuestion(id int64) (*types.Question, error)
}

type Server struct {
	service Service
	server  *http.Server
}

func NewServer(service Service) *Server {

	server := &Server{
		service: service,
	}

	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)

	// TODO: get allowed origins, headers and methods from cfg
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кеширования предзапроса (в секундах)
	}))

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)
	router.Post("/api/quizzes", server.CreateQuiz())
	router.Post("/api/quizzes/{quiz_id}/questions", server.CreateQuestion())
	router.Get("/api/questions/{quiz_id}", server.GetQuestion())

	// TODO: add timeouts
	serverCfg := &http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: router,
	}
	server.server = serverCfg

	return server
}

func (s *Server) Run() {
	slog.Info("Server started on port " + viper.GetString("port"))
	err := s.server.ListenAndServe()
	panic(err)
}
