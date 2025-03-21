package user

import (
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/services/auth"
)

type LoginHandler struct {
	authService *auth.AuthService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Token   string      `json:"token,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func NewLoginHandler(cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		authService: auth.NewAuthService(cfg),
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Falha ao processar requisição", http.StatusBadRequest)
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{
			Status:  409,
			Message: err.Error(),
		})
		return
	}

	// Criar payload similar ao da implementação original
	payload := map[string]interface{}{
		"id":        user.ID,
		"cpf":       user.CPF,
		"name":      user.Nome,
		"email":     user.Email,
		"lotacao":   user.Lotacao,
		"matricula": user.Matricula,
		"admin":     user.Admin,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Status:  201,
		Message: "Bem-vindo!",
		Token:   token,
		Payload: payload,
	})
}