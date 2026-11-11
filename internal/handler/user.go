package handler 

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"service-xp/internal/service"
	"service-xp/internal/repo"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

// Сделать один обработчик HandleUsers, внутри switch r.Method: 
// GET → GetUser, POST → CreateUser, остальное → 405.+
//  Заменить старый маршрут на новый.
// response - ответ, request - запрос
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
		case http.MethodGet:
			h.GetUser(w, r)
		case http.MethodPost:
			h.CreateUser(w, r)
		default: 
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	//1 создать структуру для приема json формата
	var req struct { 
		Name string `json:"name"`
		Email string `json:"email"`
	}
    //2 проверка json структуры
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// 3. Контекст с таймаутом(решает вопросы таймаута и отмены транзакции(нужен для любой функции ходящей в бд))
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// 4. вызвать service
	u, err := h.svc.CreateUser(ctx, req.Name, req.Email)
	if err != nil{
	switch {
		case errors.Is(err, service.ErrInvalidName):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid name"})
		case errors.Is(err, service.ErrInvalidEmail):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email"})
		default:
            writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
        }
		return 
	}

	// 5. Успех — 201
    writeJSON(w, http.StatusCreated, u)

}


func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	u, err := h.svc.GetUser(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email"})
		case errors.Is(err, repo.ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, u)
}


	func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

