package relacionamento

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
	lotacao := r.URL.Query().Get("lotacao")
	cpfCnpj := r.URL.Query().Get("cpfCnpj")
	dataInicio := r.URL.Query().Get("dataInicio")
	dataFim := r.URL.Query().Get("dataFim")
	numProcesso := r.URL.Query().Get("numProcesso")
	motivo := r.URL.Query().Get("motivo")
	caso := r.URL.Query().Get("caso")
	
	// O token já foi validado pelo middleware de autenticação

	resultado, err := h.ccsService.ConsultarRelacionamento(cpfCnpj, dataInicio, dataFim, numProcesso, motivo, cpfResponsavel, lotacao, caso)
	if err != nil {
		http.Error(w, "Erro ao consultar relacionamentos CCS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}