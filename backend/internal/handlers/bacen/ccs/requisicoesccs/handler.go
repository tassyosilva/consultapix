package requisicoesccs

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/services/bacen"
)

type Handler struct {
	ccsService *bacen.CCSService
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		ccsService: bacen.NewCCSService(cfg),
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros da URL
	cpfResponsavel := r.URL.Query().Get("cpfResponsavel")
	
	// O token já foi validado pelo middleware de autenticação

	requisicoes, err := h.ccsService.BuscarRequisicoesCCS(cpfResponsavel)
	if err != nil {
		http.Error(w, "Erro ao buscar requisições CCS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requisicoes)
}