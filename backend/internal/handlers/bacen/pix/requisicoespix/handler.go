package requisicoespix

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/repository"
)

type Handler struct {
	pixRepo *repository.PixRepository
}

func NewHandler() *Handler {
	return &Handler{
		pixRepo: repository.NewPixRepository(),
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros da URL
	cpfResponsavel := r.URL.Query().Get("cpfCnpj")
	
	// O token já foi validado pelo middleware de autenticação

	requisicoes, err := h.pixRepo.BuscarRequisicoesPix(cpfResponsavel)
	if err != nil {
		http.Error(w, "Erro ao buscar requisições PIX", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requisicoes)
}