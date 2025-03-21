package cpfcnpj

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/services/bacen"
)

type Handler struct {
	pixService *bacen.PixService
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		pixService: bacen.NewPixService(cfg),
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros da URL
	cpfResponsavel := r.URL.Query().Get("cpfResponsavel")
	lotacao := r.URL.Query().Get("lotacao")
	cpfCnpj := r.URL.Query().Get("cpfCnpj")
	motivo := r.URL.Query().Get("motivo")
	caso := r.URL.Query().Get("caso")
	
	// O token já foi validado pelo middleware de autenticação

	resultado, err := h.pixService.ConsultarPorCPFCNPJ(cpfCnpj, motivo, cpfResponsavel, lotacao, caso)
	if err != nil {
		http.Error(w, "Erro ao consultar por CPF/CNPJ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}