package detalhamento

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	numeroRequisicao := r.URL.Query().Get("numeroRequisicao")
	cpfCnpj := r.URL.Query().Get("cpfCnpj")
	cnpjResponsavel := r.URL.Query().Get("cnpjResponsavel")
	cnpjParticipante := r.URL.Query().Get("cnpjParticipante")
	dataInicioRelacionamento := r.URL.Query().Get("dataInicioRelacionamento")
	idRelacionamento := r.URL.Query().Get("idRelacionamento")
	nomeBancoResponsavel := r.URL.Query().Get("nomeBancoResponsavel")
	
	// Converter ID para inteiro
	id, err := strconv.Atoi(idRelacionamento)
	if err != nil {
		http.Error(w, "ID de relacionamento inválido", http.StatusBadRequest)
		return
	}
	
	// O token já foi validado pelo middleware de autenticação

	resultado, err := h.ccsService.SolicitarDetalhamento(numeroRequisicao, cpfCnpj, cnpjResponsavel, cnpjParticipante, dataInicioRelacionamento, nomeBancoResponsavel, id)
	if err != nil {
		http.Error(w, "Erro ao solicitar detalhamento CCS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}