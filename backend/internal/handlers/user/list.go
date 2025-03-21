package user

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/repository"
)

type ListHandler struct {
	userRepo *repository.UserRepository
}

func NewListHandler() *ListHandler {
	return &ListHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// O middleware de autenticação já validou o token, então podemos prosseguir
	usuarios, err := h.userRepo.GetAll()
	if err != nil {
		http.Error(w, "Erro ao listar usuários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuarios)
}