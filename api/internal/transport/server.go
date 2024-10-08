package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Corray333/quiz/internal/types"
	"github.com/Corray333/quiz/pkg/server/auth"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/exp/rand"
)

const MaxFileSize = 5 << 20

type Question interface {
	GetType() string
}

type Service interface {
	CreateQuiz(quiz *types.Quiz) (int64, error)
	CreateQuestion(question *types.Question) (int64, error)
	GetQuestion(id int64) (*types.Question, error)
	ListQuizzes(offset int) ([]types.Quiz, error)
	GetQuiz(id int64) (*types.Quiz, error)
	GetAnswers(userID int64, quizID int64) ([]types.Answer, error)
	GetUserAnswers(userID int64) ([]types.Answer, error)
	GetQuizAnswers(quiz_id int64, offset int) ([][]types.Answer, error)
	UpdateQuiz(quiz *types.Quiz) error
	DeleteQuiz(id int64) error

	CreateAdmin(username string) error
	IsAdminById(id int64) (bool, error)
	GetAdmins() ([]types.Admin, error)
	DeleteAdmin(id int64) error
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

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кеширования предзапроса (в секундах)
	}))

	routerAdmin := router.With(auth.NewAuthMiddleware(service))

	// TODO: get allowed origins, headers and methods from cfg

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)
	routerAdmin.Post("/api/quizzes", server.CreateQuiz())
	routerAdmin.Patch("/api/quizzes", server.UpdateQuiz())
	routerAdmin.Post("/api/quizzes/{quiz_id}/questions", server.CreateQuestion())
	routerAdmin.Get("/api/questions/{quiz_id}", server.GetQuestion())
	routerAdmin.Get("/api/quizzes", server.ListQuizzes())
	routerAdmin.Get("/api/quizzes/{quiz_id}/answers", server.GetAnswers())
	routerAdmin.Delete("/api/quizzes/{quiz_id}", server.DeleteQuiz())
	routerAdmin.Get("/api/quizzes/{quiz_id}", server.GetQuiz())
	routerAdmin.Post("/api/upload/image", server.UploadImage())
	routerAdmin.Get("/api/login", server.LogIn())

	fs := http.FileServer(http.Dir("../files/images"))
	router.Handle("/api/images/*", http.StripPrefix("/api/images", fs))

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

func (s *Server) UploadImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(MaxFileSize); err != nil {
			slog.Error("error parsing multipart form: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			slog.Error("error getting file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		rand.Seed(uint64(time.Now().UnixNano()))
		randomStr := ""
		for i := 0; i < 10; i++ {
			randomStr += strconv.Itoa(rand.Intn(10))
		}
		fileName := randomStr + filepath.Ext(header.Filename)

		_, err = os.Stat("../files/images/" + fileName)
		if err == nil {
			slog.Error("file already exists: " + fileName)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !os.IsNotExist(err) {
			slog.Error("error getting file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newFile, err := os.Create("../files/images/" + fileName)
		if err != nil {
			slog.Error("error creating file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		if _, err := io.Copy(newFile, file); err != nil {
			slog.Error("error copying file: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println()
		fmt.Println("Filename: " + os.Getenv("BASE_URL") + "/api/images/" + fileName)
		fmt.Println()

		if err := json.NewEncoder(w).Encode(struct {
			URL string `json:"url"`
		}{
			URL: os.Getenv("BASE_URL") + "/api/images/" + fileName,
		}); err != nil {
			slog.Error("error encoding or sending file name: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	}
}
