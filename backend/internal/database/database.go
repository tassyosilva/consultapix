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
	
	// Executar migrações
	if err = SetupTables(DB); err != nil {
		return fmt.Errorf("erro ao configurar tabelas: %w", err)
	}
	return nil
}

// GetDB retorna a instância do banco de dados
func GetDB() *sql.DB {
	return DB
}

// RunMigrations executa as migrações do banco de dados
func RunMigrations(db *sql.DB) error {
	// Implementação das migrações aqui
	// Exemplo simples: criar tabela se não existir
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pix_queries (
			id SERIAL PRIMARY KEY,
			key_type VARCHAR(20) NOT NULL,
			key_value VARCHAR(100) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de consultas: %w", err)
	}
	
	log.Println("Migrações executadas com sucesso")
	return nil
}