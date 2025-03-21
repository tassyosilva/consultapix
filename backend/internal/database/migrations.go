package database

import (
	"database/sql"
	"log"
)

// SetupTables verifica e cria as tabelas necessárias no banco de dados
func SetupTables(db *sql.DB) error {
	log.Println("Executando migrações do banco de dados...")

	// Criar tabela de usuários
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS usuario (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(255) NOT NULL,
			cpf VARCHAR(20) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			lotacao VARCHAR(255),
			matricula VARCHAR(50) UNIQUE,
			admin BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'usuario' verificada/criada com sucesso")

	// Tabelas para PIX
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS requisicao_pix (
			id SERIAL PRIMARY KEY,
			data TIMESTAMP NOT NULL,
			cpf_responsavel VARCHAR(20) NOT NULL,
			lotacao VARCHAR(255),
			caso VARCHAR(255),
			tipo_busca VARCHAR(50) NOT NULL,
			chave_busca VARCHAR(255) NOT NULL,
			motivo_busca VARCHAR(255) NOT NULL,
			resultado VARCHAR(255) NOT NULL,
			vinculos JSONB,
			autorizado BOOLEAN NOT NULL,
			cpf_autorizacao VARCHAR(20),
			nome_autorizacao VARCHAR(255),
			data_hora_autorizacao VARCHAR(50),
			token_autorizacao VARCHAR(255)
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'requisicao_pix' verificada/criada com sucesso")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chave_pix (
			id SERIAL PRIMARY KEY,
			chave VARCHAR(255),
			tipo_chave VARCHAR(50),
			status VARCHAR(50),
			data_abertura_reivindicacao VARCHAR(50),
			cpf_cnpj VARCHAR(20),
			nome_proprietario VARCHAR(255),
			nome_fantasia VARCHAR(255),
			participante VARCHAR(20),
			agencia VARCHAR(20),
			numero_conta VARCHAR(50),
			tipo_conta VARCHAR(50),
			data_abertura_conta VARCHAR(50),
			proprietario_da_chave_desde VARCHAR(50),
			data_criacao VARCHAR(50),
			ultima_modificacao VARCHAR(50),
			numero_banco VARCHAR(10),
			nome_banco VARCHAR(255),
			cpf_cnpj_busca VARCHAR(20),
			nome_proprietario_busca VARCHAR(255),
			id_requisicao INT NOT NULL,
			FOREIGN KEY (id_requisicao) REFERENCES requisicao_pix(id) ON DELETE CASCADE,
			UNIQUE(chave, id_requisicao)
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'chave_pix' verificada/criada com sucesso")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS evento_chave_pix (
			id SERIAL PRIMARY KEY,
			tipo_evento VARCHAR(50),
			motivo_evento VARCHAR(255),
			data_evento VARCHAR(50),
			chave VARCHAR(255),
			tipo_chave VARCHAR(50),
			cpf_cnpj VARCHAR(20),
			nome_proprietario VARCHAR(255),
			nome_fantasia VARCHAR(255),
			participante VARCHAR(20),
			agencia VARCHAR(20),
			numero_conta VARCHAR(50),
			tipo_conta VARCHAR(50),
			data_abertura_conta VARCHAR(50),
			numero_banco VARCHAR(10),
			nome_banco VARCHAR(255),
			id_chave INT NOT NULL,
			FOREIGN KEY (id_chave) REFERENCES chave_pix(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'evento_chave_pix' verificada/criada com sucesso")

	// Tabelas para CCS
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS requisicao_relacionamento_ccs (
			id SERIAL PRIMARY KEY,
			data_requisicao VARCHAR(50) NOT NULL,
			data_inicio_consulta VARCHAR(50),
			data_fim_consulta VARCHAR(50),
			cpf_cnpj_consulta VARCHAR(20),
			numero_processo VARCHAR(50),
			motivo_busca VARCHAR(255),
			cpf_responsavel VARCHAR(20),
			lotacao VARCHAR(255),
			caso VARCHAR(255),
			numero_requisicao VARCHAR(50),
			cpf_cnpj VARCHAR(20),
			tipo_pessoa VARCHAR(10),
			nome VARCHAR(255),
			autorizado BOOLEAN NOT NULL,
			cpf_autorizacao VARCHAR(20),
			nome_autorizacao VARCHAR(255),
			data_hora_autorizacao VARCHAR(50),
			token_autorizacao VARCHAR(255),
			status VARCHAR(50) NOT NULL,
			detalhamento BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'requisicao_relacionamento_ccs' verificada/criada com sucesso")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS relacionamento_ccs (
			id SERIAL PRIMARY KEY,
			numero_requisicao VARCHAR(50),
			id_pessoa VARCHAR(20),
			nome_pessoa VARCHAR(255),
			tipo_pessoa VARCHAR(10),
			cnpj_responsavel VARCHAR(20),
			numero_banco_responsavel VARCHAR(10),
			nome_banco_responsavel VARCHAR(255),
			cnpj_participante VARCHAR(20),
			numero_banco_participante VARCHAR(10),
			nome_banco_participante VARCHAR(255),
			data_inicio_relacionamento VARCHAR(50),
			data_fim_relacionamento VARCHAR(50),
			id_requisicao INT NOT NULL,
			data_requisicao_detalhamento VARCHAR(50),
			status_detalhamento VARCHAR(50) NOT NULL DEFAULT 'Nao Solicitado',
			responde_detalhamento BOOLEAN,
			resposta BOOLEAN NOT NULL DEFAULT FALSE,
			codigo_resposta VARCHAR(50),
			codigo_if_resposta VARCHAR(50),
			nuop_resposta VARCHAR(50),
			FOREIGN KEY (id_requisicao) REFERENCES requisicao_relacionamento_ccs(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'relacionamento_ccs' verificada/criada com sucesso")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bem_direito_valor_ccs (
			id SERIAL PRIMARY KEY,
			cnpj_participante VARCHAR(20),
			tipo VARCHAR(50),
			agencia VARCHAR(20),
			conta VARCHAR(50),
			vinculo VARCHAR(50),
			nome_pessoa VARCHAR(255),
			data_inicio VARCHAR(50),
			data_fim VARCHAR(50),
			id_relacionamento INT NOT NULL,
			FOREIGN KEY (id_relacionamento) REFERENCES relacionamento_ccs(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'bem_direito_valor_ccs' verificada/criada com sucesso")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vinculados_bdv_ccs (
			id SERIAL PRIMARY KEY,
			id_bdv INT NOT NULL,
			data_inicio VARCHAR(50),
			data_fim VARCHAR(50),
			id_pessoa VARCHAR(20),
			nome_pessoa VARCHAR(255),
			nome_pessoa_receita VARCHAR(255),
			tipo VARCHAR(50),
			FOREIGN KEY (id_bdv) REFERENCES bem_direito_valor_ccs(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Tabela 'vinculados_bdv_ccs' verificada/criada com sucesso")

	log.Println("Todas as migrações executadas com sucesso!")
	return nil
}