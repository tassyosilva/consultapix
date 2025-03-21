package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/routes"
	"github.com/tassyosilva/consultapix/internal/database"
)

func main() {
	// Definir variáveis de ambiente
	os.Setenv("DATABASE_URL", "postgresql://usuario:senha@localhost:5432/consultapixccs?schema=public")
	os.Setenv("usernameBC", "USUÁRIO FORNECIDO PELO BACEN. EX. eju**.s-apiccs")
	os.Setenv("passwordBC", "SENHA DE ACESSO")
	
	// Inicializar configuração
	cfg := config.NewConfig()

	// Inicializar banco de dados
	err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	// Criar router
	router := mux.NewRouter()

	// Configurar rotas
	routes.SetupRoutes(router, cfg)

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Permitir todas as origens em desenvolvimento
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not exceeded by any browser
	})

	// Encapsular router com middleware CORS
	handler := c.Handler(router)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor iniciado na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}