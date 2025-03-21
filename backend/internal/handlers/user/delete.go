package user

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/repository"
)

type DeleteHandler struct {
	userRepo *repository.UserRepository
}

type DeleteRequest struct {
	ID int `json:"id"`
}

type DeleteResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Falha ao processar requisição", http.StatusBadRequest)
		return
	}

	err := h.userRepo.Delete(req.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteResponse{
			Status:  500,
			Message: "Erro ao deletar usuário",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DeleteResponse{
		Status:  201,
		Message: "Usuário Deletado!",
	})
}