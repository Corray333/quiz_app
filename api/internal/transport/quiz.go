package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Corray333/quiz/internal/types"
	"github.com/go-chi/chi/v5"
)

type CreateQuizRequest struct {
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description,omitempty" db:"description"`
	CreatedAt   int64             `json:"createdAt" db:"created_at"`
	Cover       string            `json:"cover" db:"cover"`
	Type        string            `json:"type" db:"type"`
	Questions   []json.RawMessage `json:"questions"`
}

type CreateQuizResponse struct {
	ID int64 `json:"id"`
}

// func (s Server) CreateQuiz() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		var req CreateQuizRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		quiz := &types.Quiz{
// 			Title:       req.Title,
// 			Description: req.Description,
// 			CreatedAt:   req.CreatedAt,
// 			Cover:       req.Cover,
// 			Type:        req.Type,
// 		}
// 		id, err := s.service.CreateQuiz(quiz)
// 		if err != nil {
// 			slog.Error("create quiz error: " + err.Error())
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		if err := json.NewEncoder(w).Encode(CreateQuizResponse{
// 			ID: id,
// 		}); err != nil {
// 			slog.Error("encode response error: " + err.Error())
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

func (s Server) CreateQuiz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateQuizRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		quiz := &types.Quiz{
			Title:       req.Title,
			Description: req.Description,
			CreatedAt:   req.CreatedAt,
			Cover:       req.Cover,
			Type:        req.Type,
			Questions:   req.Questions,
		}
		id, err := s.service.CreateQuiz(quiz)
		if err != nil {
			slog.Error("create quiz error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(CreateQuizResponse{
			ID: id,
		}); err != nil {
			slog.Error("encode response error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type CreateQuestionRequest struct {
	Type     string          `json:"type" db:"type"`
	Question json.RawMessage `json:"question" db:"question"`
}

type CreateQuestionResponse struct {
	ID int64 `json:"id"`
}

func (s Server) CreateQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		quizStrID := chi.URLParam(r, "quiz_id")
		if quizStrID == "" {
			slog.Error("create question error: quiz_id url param is empty")
			http.Error(w, "quiz_id url param is empty", http.StatusBadRequest)
			return
		}
		quizID, err := strconv.ParseInt(quizStrID, 10, 64)
		if err != nil {
			slog.Error("create question error: error while parsing quiz_ad url param: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		var req CreateQuestionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Type == "" {
			slog.Error("create question error: question type is empty")
			http.Error(w, "question type is empty", http.StatusBadRequest)
			return
		}
		// TODO: add better validation

		question := &types.Question{
			Type:   req.Type,
			QuizID: quizID,
			Data:   req.Question,
		}

		// switch req.Type {
		// case types.QuestionTypeText:
		// 	question := &types.QuestionText{}
		// 	if err := json.Unmarshal(req.Question, question); err != nil {
		// 		slog.Error("create question error: unmarshalling error: " + err.Error())
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// 	}
		// 	question.QuizID = req.QuizID
		// 	question.Type = req.Type
		// case types.QuestionTypeSelect:
		// 	question := &types.QuestionSelect{}
		// 	if err := json.Unmarshal(req.Question, question); err != nil {
		// 		slog.Error("create question error: unmarshalling error: " + err.Error())
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// 	}
		// 	question.QuizID = req.QuizID
		// 	question.Type = req.Type
		// case types.QuestionTypeMultiSelect:
		// 	question := &types.QuestionMultiSelect{}
		// 	if err := json.Unmarshal(req.Question, question); err != nil {
		// 		slog.Error("create question error: unmarshalling error: " + err.Error())
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// 	}
		// 	question.QuizID = req.QuizID
		// 	question.Type = req.Type
		// }

		id, err := s.service.CreateQuestion(question)
		if err != nil {
			slog.Error("create question error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(CreateQuestionResponse{ID: id}); err != nil {
			slog.Error("encode response error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s Server) GetQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := strconv.ParseInt(chi.URLParam(r, "question_id"), 10, 64)
		if err != nil {
			slog.Error("get question error: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		question, err := s.service.GetQuestion(id)

		if err != nil {
			slog.Error("get question error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(question); err != nil {
			slog.Error("encode response error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s Server) ListQuizzes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		quizzes, err := s.service.ListQuizzes(offset)
		if err != nil {
			slog.Error("list quizzes error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(quizzes); err != nil {
			slog.Error("encode response error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s Server) GetQuiz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.ParseInt(chi.URLParam(r, "quiz_id"), 10, 64)
		if err != nil {
			slog.Error("get quiz error: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		quiz, err := s.service.GetQuiz(id)
		if err != nil {
			slog.Error("get quiz error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(quiz); err != nil {
			slog.Error("encode response error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s Server) GetAnswers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.ParseInt(chi.URLParam(r, "quiz_id"), 10, 64)
		if err != nil {
			slog.Error("get answers error: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		offsetStr := r.URL.Query().Get("offset")
		offset := 0
		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
		}
		if err != nil {
			slog.Error("get answers error: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		answers, err := s.service.GetAllAnswers(id, offset)
		if err != nil {
			slog.Error("get answers error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(answers); err != nil {
			slog.Error("encode response error: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
