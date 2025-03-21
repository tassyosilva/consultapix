package processafilaccs

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
	resultados, err := h.ccsService.ProcessarFilaCCS()
	if err != nil {
		http.Error(w, "Erro ao processar fila CCS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultados)
}