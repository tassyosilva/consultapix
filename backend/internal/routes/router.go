package routes

import (
	_"net/http"

	"github.com/gorilla/mux"
	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/handlers/bacen/ccs/detalhamento"
	"github.com/tassyosilva/consultapix/internal/handlers/bacen/ccs/relacionamento"
	"github.com/tassyosilva/consultapix/internal/handlers/bacen/ccs/requisicoesccs"
	"github.com/tassyosilva/consultapix/internal/handlers/bacen/pix/chave"
	"github.com/tassyosilva/consultapix/internal/handlers/bacen/pix/cpfcnpj"
	"github.com/tassyosilva/consultapix/internal/handlers/bacen/pix/requisicoespix"
	"github.com/tassyosilva/consultapix/internal/handlers/user"
	"github.com/tassyosilva/consultapix/internal/handlers/utils/processafilaccs"
	"github.com/tassyosilva/consultapix/internal/handlers/utils/recebebdvccs"
	"github.com/tassyosilva/consultapix/internal/middleware"
)

func SetupRoutes(router *mux.Router, cfg *config.Config) {
	// Middleware de autenticação
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Rotas públicas
	router.HandleFunc("/api/user/login", user.NewLoginHandler(cfg).Handle).Methods("POST")
	router.HandleFunc("/api/user/register", user.NewRegisterHandler().Handle).Methods("POST")

	// Rotas protegidas por autenticação
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(authMiddleware.Authenticate)

	// Rotas de usuário
	protectedRouter.HandleFunc("/user/list", user.NewListHandler().Handle).Methods("GET")
	protectedRouter.HandleFunc("/user/edit", user.NewEditHandler().Handle).Methods("POST")
	protectedRouter.HandleFunc("/user/delete", user.NewDeleteHandler().Handle).Methods("POST")

	// Rotas PIX
	protectedRouter.HandleFunc("/bacen/pix/chave", chave.NewHandler(cfg).Handle).Methods("GET")
	protectedRouter.HandleFunc("/bacen/pix/cpfCnpj", cpfcnpj.NewHandler(cfg).Handle).Methods("GET")
	protectedRouter.HandleFunc("/bacen/pix/requisicoespix", requisicoespix.NewHandler().Handle).Methods("GET")
	
	// Rotas CCS
	protectedRouter.HandleFunc("/bacen/ccs/relacionamento", relacionamento.NewHandler(cfg).Handle).Methods("GET")
	protectedRouter.HandleFunc("/bacen/ccs/detalhamento", detalhamento.NewHandler(cfg).Handle).Methods("GET")
	protectedRouter.HandleFunc("/bacen/ccs/requisicoesccs", requisicoesccs.NewHandler(cfg).Handle).Methods("GET")
	
	// Rotas para processamento em segundo plano
	router.HandleFunc("/api/utils/processaFilaCCS", processafilaccs.NewHandler(cfg).Handle).Methods("GET")
	router.HandleFunc("/api/utils/recebeBDVCCS", recebebdvccs.NewHandler(cfg).Handle).Methods("GET")
}