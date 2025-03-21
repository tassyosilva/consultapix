package repository

import (
	"database/sql"
	"encoding/json"
	_"errors"

	"github.com/tassyosilva/consultapix/internal/database"
	"github.com/tassyosilva/consultapix/internal/database/models"
)

type PixRepository struct {
	DB *sql.DB
}

func NewPixRepository() *PixRepository {
	return &PixRepository{
		DB: database.GetDB(),
	}
}

// CriarRequisicaoPix cria uma nova requisição PIX e retorna o ID
func (r *PixRepository) CriarRequisicaoPix(req *models.RequisicaoPix) (int, error) {
	// Converter vinculos para JSON
	vinculosJSON, err := json.Marshal(req.Vinculos)
	if err != nil {
		return 0, err
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Inserir requisição PIX
	insertQuery := `
		INSERT INTO requisicao_pix (
			data, cpf_responsavel, lotacao, caso, tipo_busca, chave_busca, 
			motivo_busca, resultado, vinculos, autorizado, cpf_autorizacao, 
			nome_autorizacao, data_hora_autorizacao, token_autorizacao
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	var id int
	err = tx.QueryRow(
		insertQuery,
		req.Data, req.CPFResponsavel, req.Lotacao, req.Caso, req.TipoBusca,
		req.ChaveBusca, req.MotivoBusca, req.Resultado, vinculosJSON, req.Autorizado,
		req.CPFAutorizacao, req.NomeAutorizacao, req.DataHoraAutorizacao, req.TokenAutorizacao,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Se há chaves PIX, inseri-las
	if len(req.Chaves) > 0 {
		for _, chave := range req.Chaves {
			chaveID, err := r.inserirChavePix(tx, chave, id)
			if err != nil {
				return 0, err
			}

			// Se há eventos, inseri-los
			if len(chave.EventosVinculo) > 0 {
				for _, evento := range chave.EventosVinculo {
					err = r.inserirEventoChavePix(tx, evento, chaveID)
					if err != nil {
						return 0, err
					}
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

// inserirChavePix insere uma chave PIX e retorna o ID
func (r *PixRepository) inserirChavePix(tx *sql.Tx, chave models.ChavePix, idRequisicao int) (int, error) {
	insertQuery := `
		INSERT INTO chave_pix (
			chave, tipo_chave, status, data_abertura_reivindicacao, cpf_cnpj, 
			nome_proprietario, nome_fantasia, participante, agencia, numero_conta, 
			tipo_conta, data_abertura_conta, proprietario_da_chave_desde, data_criacao, 
			ultima_modificacao, numero_banco, nome_banco, cpf_cnpj_busca, 
			nome_proprietario_busca, id_requisicao
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id
	`
	var id int
	err := tx.QueryRow(
		insertQuery,
		chave.Chave, chave.TipoChave, chave.Status, chave.DataAberturaReivindicacao,
		chave.CPFCNPJ, chave.NomeProprietario, chave.NomeFantasia, chave.Participante,
		chave.Agencia, chave.NumeroConta, chave.TipoConta, chave.DataAberturaConta,
		chave.ProprietarioDaChaveDesde, chave.DataCriacao, chave.UltimaModificacao,
		chave.NumeroBanco, chave.NomeBanco, chave.CPFCNPJBusca, chave.NomeProprietarioBusca,
		idRequisicao,
	).Scan(&id)
	return id, err
}

// inserirEventoChavePix insere um evento de chave PIX
func (r *PixRepository) inserirEventoChavePix(tx *sql.Tx, evento models.EventoChavePix, idChave int) error {
	insertQuery := `
		INSERT INTO evento_chave_pix (
			tipo_evento, motivo_evento, data_evento, chave, tipo_chave, cpf_cnpj, 
			nome_proprietario, nome_fantasia, participante, agencia, numero_conta, 
			tipo_conta, data_abertura_conta, numero_banco, nome_banco, id_chave
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := tx.Exec(
		insertQuery,
		evento.TipoEvento, evento.MotivoEvento, evento.DataEvento, evento.Chave,
		evento.TipoChave, evento.CPFCNPJ, evento.NomeProprietario, evento.NomeFantasia,
		evento.Participante, evento.Agencia, evento.NumeroConta, evento.TipoConta,
		evento.DataAberturaConta, evento.NumeroBanco, evento.NomeBanco, idChave,
	)
	return err
}

// BuscarRequisicoesPix busca todas as requisições PIX de um CPF responsável
func (r *PixRepository) BuscarRequisicoesPix(cpfResponsavel string) ([]models.RequisicaoPix, error) {
	query := `
		SELECT id, data, cpf_responsavel, lotacao, caso, tipo_busca, chave_busca, 
			motivo_busca, resultado, vinculos, autorizado, cpf_autorizacao, 
			nome_autorizacao, data_hora_autorizacao, token_autorizacao
		FROM requisicao_pix
		WHERE cpf_responsavel = $1
		ORDER BY id DESC
	`
	rows, err := r.DB.Query(query, cpfResponsavel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requisicoes []models.RequisicaoPix
	for rows.Next() {
		var req models.RequisicaoPix
		var vinculosJSON []byte
		
		err := rows.Scan(
			&req.ID, &req.Data, &req.CPFResponsavel, &req.Lotacao, &req.Caso, 
			&req.TipoBusca, &req.ChaveBusca, &req.MotivoBusca, &req.Resultado, 
			&vinculosJSON, &req.Autorizado, &req.CPFAutorizacao, &req.NomeAutorizacao, 
			&req.DataHoraAutorizacao, &req.TokenAutorizacao,
		)
		if err != nil {
			return nil, err
		}

		// Converter JSON para interface{}
		if len(vinculosJSON) > 0 {
			err = json.Unmarshal(vinculosJSON, &req.Vinculos)
			if err != nil {
				return nil, err
			}
		}

		requisicoes = append(requisicoes, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requisicoes, nil
}