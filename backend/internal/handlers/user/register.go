package user

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/database/models"
	"github.com/tassyosilva/consultapix/internal/repository"
)

type RegisterHandler struct {
	userRepo *repository.UserRepository
}

type RegisterRequest struct {
	Nome      string `json:"nome"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Lotacao   string `json:"lotacao"`
	Matricula string `json:"matricula"`
	Admin     bool   `json:"admin"`
}

type RegisterResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Falha ao processar requisição", http.StatusBadRequest)
		return
	}

	user := &models.Usuario{
		Nome:      req.Nome,
		CPF:       req.CPF,
		Email:     req.Email,
		Password:  req.Password,
		Lotacao:   req.Lotacao,
		Matricula: req.Matricula,
		Admin:     req.Admin,
	}

	_, err := h.userRepo.Create(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RegisterResponse{
			Status:  409,
			Message: "Usuário já cadastrado",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RegisterResponse{
		Status:  201,
		Message: "Obrigado pelo cadastro!",
	})
}