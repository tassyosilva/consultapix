package user

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/database/models"
	"github.com/tassyosilva/consultapix/internal/repository"
)

type EditHandler struct {
	userRepo *repository.UserRepository
}

type EditRequest struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Lotacao   string `json:"lotacao"`
	Matricula string `json:"matricula"`
	Admin     bool   `json:"admin"`
}

type EditResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewEditHandler() *EditHandler {
	return &EditHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *EditHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req EditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Falha ao processar requisição", http.StatusBadRequest)
		return
	}

	user := &models.Usuario{
		ID:        req.ID,
		Nome:      req.Nome,
		CPF:       req.CPF,
		Email:     req.Email,
		Password:  req.Password,
		Lotacao:   req.Lotacao,
		Matricula: req.Matricula,
		Admin:     req.Admin,
	}

	err := h.userRepo.Update(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(EditResponse{
			Status:  500,
			Message: "Erro ao atualizar usuário",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(EditResponse{
		Status:  201,
		Message: "Usuário Atualizado!",
	})
}