package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/tassyosilva/consultapix/internal/database"
	"github.com/tassyosilva/consultapix/internal/database/models"
)

type CCSRepository struct {
	DB *sql.DB
}

func NewCCSRepository() *CCSRepository {
	return &CCSRepository{
		DB: database.GetDB(),
	}
}

// CriarRequisicaoRelacionamentoCCS cria uma nova requisição de relacionamento CCS
func (r *CCSRepository) CriarRequisicaoRelacionamentoCCS(req *models.RequisicaoRelacionamentoCCS) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Inserir requisição de relacionamento CCS
	insertQuery := `
		INSERT INTO requisicao_relacionamento_ccs (
			data_requisicao, data_inicio_consulta, data_fim_consulta, cpf_cnpj_consulta,
			numero_processo, motivo_busca, cpf_responsavel, lotacao, caso,
			numero_requisicao, cpf_cnpj, tipo_pessoa, nome, autorizado,
			cpf_autorizacao, nome_autorizacao, data_hora_autorizacao, token_autorizacao,
			status, detalhamento
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id
	`
	var id int
	err = tx.QueryRow(
		insertQuery,
		req.DataRequisicao, req.DataInicioConsulta, req.DataFimConsulta, req.CPFCNPJConsulta,
		req.NumeroProcesso, req.MotivoBusca, req.CPFResponsavel, req.Lotacao, req.Caso,
		req.NumeroRequisicao, req.CPFCNPJ, req.TipoPessoa, req.Nome, req.Autorizado,
		req.CPFAutorizacao, req.NomeAutorizacao, req.DataHoraAutorizacao, req.TokenAutorizacao,
		req.Status, req.Detalhamento,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Se há relacionamentos, inseri-los
	if len(req.RelacionamentosCCS) > 0 {
		for _, relacionamento := range req.RelacionamentosCCS {
			relID, err := r.inserirRelacionamentoCCS(tx, relacionamento, id)
			if err != nil {
				return 0, err
			}

			// Se há bens/direitos/valores, inseri-los
			if len(relacionamento.BemDireitoValorCCS) > 0 {
				for _, bdv := range relacionamento.BemDireitoValorCCS {
					bdvID, err := r.inserirBemDireitoValorCCS(tx, bdv, relID)
					if err != nil {
						return 0, err
					}

					// Se há vinculados, inseri-los
					if len(bdv.Vinculados) > 0 {
						for _, vinculado := range bdv.Vinculados {
							err = r.inserirVinculadosBDVCCS(tx, vinculado, bdvID)
							if err != nil {
								return 0, err
							}
						}
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

// inserirRelacionamentoCCS insere um relacionamento CCS e retorna o ID
func (r *CCSRepository) inserirRelacionamentoCCS(tx *sql.Tx, rel models.RelacionamentoCCS, idRequisicao int) (int, error) {
	insertQuery := `
		INSERT INTO relacionamento_ccs (
			numero_requisicao, id_pessoa, nome_pessoa, tipo_pessoa, cnpj_responsavel,
			numero_banco_responsavel, nome_banco_responsavel, cnpj_participante,
			numero_banco_participante, nome_banco_participante, data_inicio_relacionamento,
			data_fim_relacionamento, id_requisicao, data_requisicao_detalhamento,
			status_detalhamento, responde_detalhamento, resposta, codigo_resposta,
			codigo_if_resposta, nuop_resposta
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id
	`
	var id int
	err := tx.QueryRow(
		insertQuery,
		rel.NumeroRequisicao, rel.IDPessoa, rel.NomePessoa, rel.TipoPessoa, rel.CNPJResponsavel,
		rel.NumeroBancoResponsavel, rel.NomeBancoResponsavel, rel.CNPJParticipante,
		rel.NumeroBancoParticipante, rel.NomeBancoParticipante, rel.DataInicioRelacionamento,
		rel.DataFimRelacionamento, idRequisicao, rel.DataRequisicaoDetalhamento,
		rel.StatusDetalhamento, rel.RespondeDetalhamento, rel.Resposta, rel.CodigoResposta,
		rel.CodigoIfResposta, rel.NuopResposta,
	).Scan(&id)
	return id, err
}

// inserirBemDireitoValorCCS insere um BDV CCS e retorna o ID
func (r *CCSRepository) inserirBemDireitoValorCCS(tx *sql.Tx, bdv models.BemDireitoValorCCS, idRelacionamento int) (int, error) {
	insertQuery := `
		INSERT INTO bem_direito_valor_ccs (
			cnpj_participante, tipo, agencia, conta, vinculo, nome_pessoa, 
			data_inicio, data_fim, id_relacionamento
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	var id int
	err := tx.QueryRow(
		insertQuery,
		bdv.CNPJParticipante, bdv.Tipo, bdv.Agencia, bdv.Conta, bdv.Vinculo,
		bdv.NomePessoa, bdv.DataInicio, bdv.DataFim, idRelacionamento,
	).Scan(&id)
	return id, err
}

// inserirVinculadosBDVCCS insere um vinculado BDV CCS
func (r *CCSRepository) inserirVinculadosBDVCCS(tx *sql.Tx, vinc models.VinculadosBDVCCS, idBDV int) error {
	insertQuery := `
		INSERT INTO vinculados_bdv_ccs (
			id_bdv, data_inicio, data_fim, id_pessoa, nome_pessoa, nome_pessoa_receita, tipo
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(
		insertQuery,
		idBDV, vinc.DataInicio, vinc.DataFim, vinc.IDPessoa, vinc.NomePessoa,
		vinc.NomePessoaReceita, vinc.Tipo,
	)
	return err
}

// BuscarRequisicoesRelacionamentoCCS busca todas as requisições CCS de um CPF responsável
func (r *CCSRepository) BuscarRequisicoesRelacionamentoCCS(cpfResponsavel string) ([]models.RequisicaoRelacionamentoCCS, error) {
	query := `
		SELECT id, data_requisicao, data_inicio_consulta, data_fim_consulta, cpf_cnpj_consulta,
			numero_processo, motivo_busca, cpf_responsavel, lotacao, caso,
			numero_requisicao, cpf_cnpj, tipo_pessoa, nome, autorizado,
			cpf_autorizacao, nome_autorizacao, data_hora_autorizacao, token_autorizacao,
			status, detalhamento
		FROM requisicao_relacionamento_ccs
		WHERE cpf_responsavel = $1
		ORDER BY id DESC
	`
	rows, err := r.DB.Query(query, cpfResponsavel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requisicoes []models.RequisicaoRelacionamentoCCS
	for rows.Next() {
		var req models.RequisicaoRelacionamentoCCS
		
		err := rows.Scan(
			&req.ID, &req.DataRequisicao, &req.DataInicioConsulta, &req.DataFimConsulta, 
			&req.CPFCNPJConsulta, &req.NumeroProcesso, &req.MotivoBusca, &req.CPFResponsavel, 
			&req.Lotacao, &req.Caso, &req.NumeroRequisicao, &req.CPFCNPJ, &req.TipoPessoa, 
			&req.Nome, &req.Autorizado, &req.CPFAutorizacao, &req.NomeAutorizacao, 
			&req.DataHoraAutorizacao, &req.TokenAutorizacao, &req.Status, &req.Detalhamento,
		)
		if err != nil {
			return nil, err
		}

		// Buscar relacionamentos
		relacionamentos, err := r.buscarRelacionamentosPorRequisicao(req.ID)
		if err != nil {
			return nil, err
		}
		req.RelacionamentosCCS = relacionamentos

		requisicoes = append(requisicoes, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requisicoes, nil
}

// buscarRelacionamentosPorRequisicao busca todos os relacionamentos CCS de uma requisição
func (r *CCSRepository) buscarRelacionamentosPorRequisicao(idRequisicao int) ([]models.RelacionamentoCCS, error) {
	query := `
		SELECT id, numero_requisicao, id_pessoa, nome_pessoa, tipo_pessoa, cnpj_responsavel,
			numero_banco_responsavel, nome_banco_responsavel, cnpj_participante,
			numero_banco_participante, nome_banco_participante, data_inicio_relacionamento,
			data_fim_relacionamento, id_requisicao, data_requisicao_detalhamento,
			status_detalhamento, responde_detalhamento, resposta, codigo_resposta,
			codigo_if_resposta, nuop_resposta
		FROM relacionamento_ccs
		WHERE id_requisicao = $1
		ORDER BY numero_banco_responsavel ASC
	`
	rows, err := r.DB.Query(query, idRequisicao)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relacionamentos []models.RelacionamentoCCS
	for rows.Next() {
		var rel models.RelacionamentoCCS
		
		err := rows.Scan(
			&rel.ID, &rel.NumeroRequisicao, &rel.IDPessoa, &rel.NomePessoa, &rel.TipoPessoa, 
			&rel.CNPJResponsavel, &rel.NumeroBancoResponsavel, &rel.NomeBancoResponsavel, 
			&rel.CNPJParticipante, &rel.NumeroBancoParticipante, &rel.NomeBancoParticipante, 
			&rel.DataInicioRelacionamento, &rel.DataFimRelacionamento, &rel.IDRequisicao, 
			&rel.DataRequisicaoDetalhamento, &rel.StatusDetalhamento, &rel.RespondeDetalhamento, 
			&rel.Resposta, &rel.CodigoResposta, &rel.CodigoIfResposta, &rel.NuopResposta,
		)
		if err != nil {
			return nil, err
		}

		// Buscar BDVs
		bdvs, err := r.buscarBDVPorRelacionamento(rel.ID)
		if err != nil {
			return nil, err
		}
		rel.BemDireitoValorCCS = bdvs

		relacionamentos = append(relacionamentos, rel)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return relacionamentos, nil
}

// buscarBDVPorRelacionamento busca todos os BDVs de um relacionamento
func (r *CCSRepository) buscarBDVPorRelacionamento(idRelacionamento int) ([]models.BemDireitoValorCCS, error) {
	query := `
		SELECT id, cnpj_participante, tipo, agencia, conta, vinculo, nome_pessoa, 
			data_inicio, data_fim, id_relacionamento
		FROM bem_direito_valor_ccs
		WHERE id_relacionamento = $1
		ORDER BY data_inicio DESC
	`
	rows, err := r.DB.Query(query, idRelacionamento)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bdvs []models.BemDireitoValorCCS
	for rows.Next() {
		var bdv models.BemDireitoValorCCS
		
		err := rows.Scan(
			&bdv.ID, &bdv.CNPJParticipante, &bdv.Tipo, &bdv.Agencia, &bdv.Conta, 
			&bdv.Vinculo, &bdv.NomePessoa, &bdv.DataInicio, &bdv.DataFim, &bdv.IDRelacionamento,
		)
		if err != nil {
			return nil, err
		}

		// Buscar vinculados
		vinculados, err := r.buscarVinculadosPorBDV(bdv.ID)
		if err != nil {
			return nil, err
		}
		bdv.Vinculados = vinculados

		bdvs = append(bdvs, bdv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bdvs, nil
}

// buscarVinculadosPorBDV busca todos os vinculados de um BDV
func (r *CCSRepository) buscarVinculadosPorBDV(idBDV int) ([]models.VinculadosBDVCCS, error) {
	query := `
		SELECT id, id_bdv, data_inicio, data_fim, id_pessoa, nome_pessoa, nome_pessoa_receita, tipo
		FROM vinculados_bdv_ccs
		WHERE id_bdv = $1
	`
	rows, err := r.DB.Query(query, idBDV)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vinculados []models.VinculadosBDVCCS
	for rows.Next() {
		var vinc models.VinculadosBDVCCS
		
		err := rows.Scan(
			&vinc.ID, &vinc.IDBDV, &vinc.DataInicio, &vinc.DataFim, &vinc.IDPessoa, 
			&vinc.NomePessoa, &vinc.NomePessoaReceita, &vinc.Tipo,
		)
		if err != nil {
			return nil, err
		}

		vinculados = append(vinculados, vinc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vinculados, nil
}

// AtualizarStatusDetalhamentoCCS atualiza o status de detalhamento de um relacionamento CCS
func (r *CCSRepository) AtualizarStatusDetalhamentoCCS(id int, status string, respondeDetalhamento, resposta bool, dataRequisicaoDetalhamento, codigoResposta, codigoIfResposta, nuopResposta string) error {
	query := `
		UPDATE relacionamento_ccs
		SET status_detalhamento = $1,
			responde_detalhamento = $2,
			resposta = $3,
			data_requisicao_detalhamento = $4,
			codigo_resposta = $5,
			codigo_if_resposta = $6,
			nuop_resposta = $7
		WHERE id = $8
	`
	_, err := r.DB.Exec(
		query, 
		status, 
		respondeDetalhamento, 
		resposta, 
		dataRequisicaoDetalhamento, 
		codigoResposta, 
		codigoIfResposta, 
		nuopResposta, 
		id,
	)
	return err
}

// BuscarRelacionamentosNaFila busca todos os relacionamentos CCS com status "Na fila"
func (r *CCSRepository) BuscarRelacionamentosNaFila() ([]models.RequisicaoRelacionamentoCCS, error) {
	query := `
		SELECT r.id, r.data_requisicao, r.data_inicio_consulta, r.data_fim_consulta, r.cpf_cnpj_consulta,
			r.numero_processo, r.motivo_busca, r.cpf_responsavel, r.lotacao, r.caso,
			r.numero_requisicao, r.cpf_cnpj, r.tipo_pessoa, r.nome, r.autorizado,
			r.cpf_autorizacao, r.nome_autorizacao, r.data_hora_autorizacao, r.token_autorizacao,
			r.status, r.detalhamento
		FROM requisicao_relacionamento_ccs r
		INNER JOIN relacionamento_ccs rc ON r.id = rc.id_requisicao
		WHERE rc.status_detalhamento = 'Na fila'
		GROUP BY r.id
		ORDER BY r.id DESC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requisicoes []models.RequisicaoRelacionamentoCCS
	for rows.Next() {
		var req models.RequisicaoRelacionamentoCCS
		
		err := rows.Scan(
			&req.ID, &req.DataRequisicao, &req.DataInicioConsulta, &req.DataFimConsulta, 
			&req.CPFCNPJConsulta, &req.NumeroProcesso, &req.MotivoBusca, &req.CPFResponsavel, 
			&req.Lotacao, &req.Caso, &req.NumeroRequisicao, &req.CPFCNPJ, &req.TipoPessoa, 
			&req.Nome, &req.Autorizado, &req.CPFAutorizacao, &req.NomeAutorizacao, 
			&req.DataHoraAutorizacao, &req.TokenAutorizacao, &req.Status, &req.Detalhamento,
		)
		if err != nil {
			return nil, err
		}

		// Buscar apenas relacionamentos na fila
		query := `
			SELECT id, numero_requisicao, id_pessoa, nome_pessoa, tipo_pessoa, cnpj_responsavel,
				numero_banco_responsavel, nome_banco_responsavel, cnpj_participante,
				numero_banco_participante, nome_banco_participante, data_inicio_relacionamento,
				data_fim_relacionamento, id_requisicao, data_requisicao_detalhamento,
				status_detalhamento, responde_detalhamento, resposta, codigo_resposta,
				codigo_if_resposta, nuop_resposta
			FROM relacionamento_ccs
			WHERE id_requisicao = $1 AND status_detalhamento = 'Na fila'
		`
		relRows, err := r.DB.Query(query, req.ID)
		if err != nil {
			return nil, err
		}

		var relacionamentos []models.RelacionamentoCCS
		for relRows.Next() {
			var rel models.RelacionamentoCCS
			
			err := relRows.Scan(
				&rel.ID, &rel.NumeroRequisicao, &rel.IDPessoa, &rel.NomePessoa, &rel.TipoPessoa, 
				&rel.CNPJResponsavel, &rel.NumeroBancoResponsavel, &rel.NomeBancoResponsavel, 
				&rel.CNPJParticipante, &rel.NumeroBancoParticipante, &rel.NomeBancoParticipante, 
				&rel.DataInicioRelacionamento, &rel.DataFimRelacionamento, &rel.IDRequisicao, 
				&rel.DataRequisicaoDetalhamento, &rel.StatusDetalhamento, &rel.RespondeDetalhamento, 
				&rel.Resposta, &rel.CodigoResposta, &rel.CodigoIfResposta, &rel.NuopResposta,
			)
			if err != nil {
				relRows.Close()
				return nil, err
			}

			relacionamentos = append(relacionamentos, rel)
		}
		relRows.Close()

		if err = relRows.Err(); err != nil {
			return nil, err
		}

		req.RelacionamentosCCS = relacionamentos
		requisicoes = append(requisicoes, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requisicoes, nil
}

// BuscarRelacionamentosAguardandoResposta busca todos os relacionamentos CCS com status "Solicitado. Aguardando..."
func (r *CCSRepository) BuscarRelacionamentosAguardandoResposta() ([]models.RequisicaoRelacionamentoCCS, error) {
	query := `
		SELECT r.id, r.data_requisicao, r.data_inicio_consulta, r.data_fim_consulta, r.cpf_cnpj_consulta,
			r.numero_processo, r.motivo_busca, r.cpf_responsavel, r.lotacao, r.caso,
			r.numero_requisicao, r.cpf_cnpj, r.tipo_pessoa, r.nome, r.autorizado,
			r.cpf_autorizacao, r.nome_autorizacao, r.data_hora_autorizacao, r.token_autorizacao,
			r.status, r.detalhamento
		FROM requisicao_relacionamento_ccs r
		INNER JOIN relacionamento_ccs rc ON r.id = rc.id_requisicao
		WHERE rc.status_detalhamento = 'Solicitado. Aguardando...'
		GROUP BY r.id
		ORDER BY r.id DESC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requisicoes []models.RequisicaoRelacionamentoCCS
	for rows.Next() {
		var req models.RequisicaoRelacionamentoCCS
		
		err := rows.Scan(
			&req.ID, &req.DataRequisicao, &req.DataInicioConsulta, &req.DataFimConsulta, 
			&req.CPFCNPJConsulta, &req.NumeroProcesso, &req.MotivoBusca, &req.CPFResponsavel, 
			&req.Lotacao, &req.Caso, &req.NumeroRequisicao, &req.CPFCNPJ, &req.TipoPessoa, 
			&req.Nome, &req.Autorizado, &req.CPFAutorizacao, &req.NomeAutorizacao, 
			&req.DataHoraAutorizacao, &req.TokenAutorizacao, &req.Status, &req.Detalhamento,
		)
		if err != nil {
			return nil, err
		}

		// Buscar apenas relacionamentos aguardando resposta
		query := `
			SELECT id, numero_requisicao, id_pessoa, nome_pessoa, tipo_pessoa, cnpj_responsavel,
				numero_banco_responsavel, nome_banco_responsavel, cnpj_participante,
				numero_banco_participante, nome_banco_participante, data_inicio_relacionamento,
				data_fim_relacionamento, id_requisicao, data_requisicao_detalhamento,
				status_detalhamento, responde_detalhamento, resposta, codigo_resposta,
				codigo_if_resposta, nuop_resposta
			FROM relacionamento_ccs
			WHERE id_requisicao = $1 AND status_detalhamento = 'Solicitado. Aguardando...'
		`
		relRows, err := r.DB.Query(query, req.ID)
		if err != nil {
			return nil, err
		}

		var relacionamentos []models.RelacionamentoCCS
		for relRows.Next() {
			var rel models.RelacionamentoCCS
			
			err := relRows.Scan(
				&rel.ID, &rel.NumeroRequisicao, &rel.IDPessoa, &rel.NomePessoa, &rel.TipoPessoa, 
				&rel.CNPJResponsavel, &rel.NumeroBancoResponsavel, &rel.NomeBancoResponsavel, 
				&rel.CNPJParticipante, &rel.NumeroBancoParticipante, &rel.NomeBancoParticipante, 
				&rel.DataInicioRelacionamento, &rel.DataFimRelacionamento, &rel.IDRequisicao, 
				&rel.DataRequisicaoDetalhamento, &rel.StatusDetalhamento, &rel.RespondeDetalhamento, 
				&rel.Resposta, &rel.CodigoResposta, &rel.CodigoIfResposta, &rel.NuopResposta,
			)
			if err != nil {
				relRows.Close()
				return nil, err
			}

			relacionamentos = append(relacionamentos, rel)
		}
		relRows.Close()

		if err = relRows.Err(); err != nil {
			return nil, err
		}

		req.RelacionamentosCCS = relacionamentos
		requisicoes = append(requisicoes, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requisicoes, nil
}