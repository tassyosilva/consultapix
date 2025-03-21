// internal/database/database.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver PostgreSQL
	"github.com/tassyosilva/consultapix/internal/config"
)

// DB é uma instância global do banco de dados
var DB *sql.DB

// Initialize inicializa a conexão com o banco de dados
func Initialize(cfg *config.Config) error {
	var err error

	// Abrir conexão com o PostgreSQL
	DB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	// Verificar se a conexão está funcionando
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("erro ao verificar conexão com banco de dados: %w", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return nil
}

// GetDB retorna a instância do banco de dados
func GetDB() *sql.DB {
	return DB
}