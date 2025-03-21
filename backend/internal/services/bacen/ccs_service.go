package bacen

import (
	_"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	_"strings"
	"time"

	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/database/models"
	"github.com/tassyosilva/consultapix/internal/repository"
)

type CCSService struct {
	config   *config.Config
	ccsRepo  *repository.CCSRepository
	client   *http.Client
}

// Estruturas para trabalhar com XML do BACEN

type RequisicaoRelacionamentoXML struct {
	XMLName         xml.Name `xml:"requisicaoRelacionamento"`
	NumeroRequisicao string   `xml:"numeroRequisicao"`
	DataMovimento    string   `xml:"dataMovimento"`
	NumeroProcesso   string   `xml:"numeroProcesso"`
	Motivo           string   `xml:"motivo"`
	Clientes         struct {
		Clientes []struct {
			ID          string `xml:"id"`
			Nome        string `xml:"nome"`
			TipoPessoa  string `xml:"tipoPessoa"`
			Relacionamentos struct {
				Relacionamentos []RelacionamentoXML `xml:"relacionamentos"`
			} `xml:"relacionamentos"`
		} `xml:"clientes"`
	} `xml:"clientes"`
}

type RelacionamentoXML struct {
	CNPJ              string `xml:"cnpj"`
	CNPJParticipante  string `xml:"cnpjParticipante"`
	ResponsavelAtivo  string `xml:"responsavelAtivo"`
	Periodos          struct {
		Periodos []struct {
			DataInicio  string `xml:"dataInicio"`
			DataFim     string `xml:"dataFim,omitempty"`
		} `xml:"periodos"`
	} `xml:"periodos"`
}

type RequisicaoDetalhamentosXML struct {
	XMLName              xml.Name `xml:"requisicaoDetalhamentos"`
	RequisicaoDetalhamento []struct {
		DataHoraRequisicao string `xml:"dataHoraRequisicao"`
		// outros campos relevantes
	} `xml:"requisicaoDetalhamento"`
}

type RespostaDetalhamentosXML struct {
	XMLName              xml.Name `xml:"respostaDetalhamentos"`
	RespostaDetalhamento []struct {
		Codigo     string `xml:"codigo"`
		CodigoIf   string `xml:"codigoIf,omitempty"`
		Nuop       string `xml:"nuop,omitempty"`
		// outros campos relevantes
	} `xml:"respostaDetalhamento"`
}

type BemDireitoValorsXML struct {
	XMLName          xml.Name `xml:"bemDireitoValors"`
	BemDireitoValor  []struct {
		CNPJParticipante  string `xml:"cnpjParticipante,omitempty"`
		Tipo              string `xml:"tipo,omitempty"`
		Agencia           string `xml:"agencia,omitempty"`
		Conta             string `xml:"conta,omitempty"`
		Vinculo           string `xml:"vinculo,omitempty"`
		NomePessoa        string `xml:"nomePessoa,omitempty"`
		DataInicio        string `xml:"dataInicio,omitempty"`
		DataFim           string `xml:"dataFim,omitempty"`
		Vinculados        struct {
			Vinculados []struct {
				IDPessoa         string `xml:"idPessoa"`
				DataInicio       string `xml:"dataInicio"`
				DataFim          string `xml:"dataFim,omitempty"`
				NomePessoa       string `xml:"nomePessoa"`
				NomePessoaReceita string `xml:"nomePessoaReceita"`
				Tipo             string `xml:"tipo"`
			} `xml:"vinculados,omitempty"`
		} `xml:"vinculados"`
	} `xml:"bemDireitoValor"`
}

func NewCCSService(cfg *config.Config) *CCSService {
	return &CCSService{
		config:  cfg,
		ccsRepo: repository.NewCCSRepository(),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// getBasicAuth retorna o cabeçalho de autenticação básica
func (s *CCSService) getBasicAuth() string {
	credentials := fmt.Sprintf("%s:%s", s.config.BacenUsername, s.config.BacenPassword)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}

// ConsultarRelacionamento consulta relacionamentos CCS de um CPF/CNPJ
func (s *CCSService) ConsultarRelacionamento(cpfCnpj, dataInicio, dataFim, numProcesso, motivo string, cpfResponsavel, lotacao, caso string) ([]models.RequisicaoRelacionamentoCCS, error) {
	url := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/requisitar-relacionamentos?id-cliente=%s&data-inicio=%s&data-fim=%s&numero-processo=%s&motivo=%s", 
		cpfCnpj, dataInicio, dataFim, numProcesso, motivo)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", s.getBasicAuth())
	
	resp, err := s.client.Do(req)
	if err != nil {
		// Registrar falha na requisição
		requisicao := &models.RequisicaoRelacionamentoCCS{
			DataRequisicao:     time.Now().Format(time.RFC3339),
			DataInicioConsulta: dataInicio,
			DataFimConsulta:    dataFim,
			CPFCNPJConsulta:    cpfCnpj,
			NumeroProcesso:     numProcesso,
			MotivoBusca:        motivo,
			CPFResponsavel:     cpfResponsavel,
			Lotacao:            lotacao,
			Caso:               caso,
			NumeroRequisicao:   "",
			CPFCNPJ:            "",
			TipoPessoa:         "",
			Nome:               "",
			Autorizado:         true,
			Status:             "Falha",
		}
		
		_, _ = s.ccsRepo.CriarRequisicaoRelacionamentoCCS(requisicao)
		
		return []models.RequisicaoRelacionamentoCCS{{
			CPFCNPJConsulta: cpfCnpj,
			Status:          "Falha",
			Nome:            "CPF ou CNPJ incorreto",
		}}, nil
	}
	defer resp.Body.Close()
	
	// Ler resposta XML
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Converter XML para estrutura
	var requisicaoXML RequisicaoRelacionamentoXML
	err = xml.Unmarshal(bodyBytes, &requisicaoXML)
	if err != nil {
		return nil, err
	}
	
	// Verificar se tem cliente e relacionamentos
	if len(requisicaoXML.Clientes.Clientes) == 0 {
		// Registrar requisição sem relacionamentos
		requisicao := &models.RequisicaoRelacionamentoCCS{
			DataRequisicao:     requisicaoXML.DataMovimento,
			DataInicioConsulta: dataInicio,
			DataFimConsulta:    dataFim,
			CPFCNPJConsulta:    cpfCnpj,
			NumeroProcesso:     requisicaoXML.NumeroProcesso,
			MotivoBusca:        requisicaoXML.Motivo,
			CPFResponsavel:     cpfResponsavel,
			Lotacao:            lotacao,
			Caso:               caso,
			NumeroRequisicao:   requisicaoXML.NumeroRequisicao,
			CPFCNPJ:            "",
			TipoPessoa:         "",
			Nome:               "",
			Autorizado:         true,
			Status:             "Sucesso",
		}
		
		_, err = s.ccsRepo.CriarRequisicaoRelacionamentoCCS(requisicao)
		if err != nil {
			return nil, err
		}
		
		return []models.RequisicaoRelacionamentoCCS{{
			CPFCNPJConsulta: cpfCnpj,
			Status:          "Sucesso",
			Nome:            "CPF / CNPJ Não possui relacionamentos no período informado",
		}}, nil
	}
	
	// Processar cliente e relacionamentos
	cliente := requisicaoXML.Clientes.Clientes[0]
	
	// Verificar se tem relacionamentos
	if len(cliente.Relacionamentos.Relacionamentos) == 0 {
		// Registrar requisição sem relacionamentos
		requisicao := &models.RequisicaoRelacionamentoCCS{
			DataRequisicao:     requisicaoXML.DataMovimento,
			DataInicioConsulta: dataInicio,
			DataFimConsulta:    dataFim,
			CPFCNPJConsulta:    cpfCnpj,
			NumeroProcesso:     requisicaoXML.NumeroProcesso,
			MotivoBusca:        requisicaoXML.Motivo,
			CPFResponsavel:     cpfResponsavel,
			Lotacao:            lotacao,
			Caso:               caso,
			NumeroRequisicao:   requisicaoXML.NumeroRequisicao,
			CPFCNPJ:            cliente.ID,
			TipoPessoa:         cliente.TipoPessoa,
			Nome:               cliente.Nome,
			Autorizado:         true,
			Status:             "Sucesso",
		}
		
		_, err = s.ccsRepo.CriarRequisicaoRelacionamentoCCS(requisicao)
		if err != nil {
			return nil, err
		}
		
		return []models.RequisicaoRelacionamentoCCS{{
			CPFCNPJConsulta: cpfCnpj,
			Status:          "Sucesso",
			Nome:            "CPF / CNPJ Não possui relacionamentos no período informado",
		}}, nil
	}
	
	// Processar relacionamentos
	relacionamentos := make([]models.RelacionamentoCCS, 0)
	
	for _, relXML := range cliente.Relacionamentos.Relacionamentos {
		// Buscar dados do participante responsável
		participanteResp, err := s.ConsultarParticipante(relXML.CNPJ)
		var numeroBancoResponsavel, nomeBancoResponsavel string
		
		if err != nil {
			numeroBancoResponsavel = "000"
			nomeBancoResponsavel = "BANCO NÃO INFORMADO"
		} else {
			numeroBancoResponsavel = fmt.Sprintf("%03d", participanteResp.CodigoCompensacao)
			nomeBancoResponsavel = participanteResp.Nome
		}
		
		// Buscar dados do participante banco
		participanteBanco, err := s.ConsultarParticipante(relXML.CNPJParticipante)
		var numeroBancoParticipante, nomeBancoParticipante string
		
		if err != nil {
			numeroBancoParticipante = "000"
			nomeBancoParticipante = "BANCO NÃO INFORMADO"
		} else {
			numeroBancoParticipante = fmt.Sprintf("%03d", participanteBanco.CodigoCompensacao)
			nomeBancoParticipante = participanteBanco.Nome
		}
		
		// Obter datas de início e fim do relacionamento
		var dataInicioRel, dataFimRel string
		if len(relXML.Periodos.Periodos) > 0 {
			dataInicioRel = relXML.Periodos.Periodos[0].DataInicio
			dataFimRel = relXML.Periodos.Periodos[0].DataFim
		}
		
		// Criar modelo de relacionamento
		relacionamento := models.RelacionamentoCCS{
			NumeroRequisicao:        requisicaoXML.NumeroRequisicao,
			IDPessoa:                cliente.ID,
			NomePessoa:              cliente.Nome,
			TipoPessoa:              cliente.TipoPessoa,
			CNPJResponsavel:         relXML.CNPJ,
			NumeroBancoResponsavel:  numeroBancoResponsavel,
			NomeBancoResponsavel:    nomeBancoResponsavel,
			CNPJParticipante:        relXML.CNPJParticipante,
			NumeroBancoParticipante: numeroBancoParticipante,
			NomeBancoParticipante:   nomeBancoParticipante,
			DataInicioRelacionamento: dataInicioRel,
			DataFimRelacionamento:   dataFimRel,
			StatusDetalhamento:      "Nao Solicitado",
		}
		
		relacionamentos = append(relacionamentos, relacionamento)
	}
	
	// Criar requisição para salvar no banco
	requisicao := &models.RequisicaoRelacionamentoCCS{
		DataRequisicao:     requisicaoXML.DataMovimento,
		DataInicioConsulta: dataInicio,
		DataFimConsulta:    dataFim,
		CPFCNPJConsulta:    cpfCnpj,
		NumeroProcesso:     requisicaoXML.NumeroProcesso,
		MotivoBusca:        requisicaoXML.Motivo,
		CPFResponsavel:     cpfResponsavel,
		Lotacao:            lotacao,
		Caso:               caso,
		NumeroRequisicao:   requisicaoXML.NumeroRequisicao,
		CPFCNPJ:            cliente.ID,
		TipoPessoa:         cliente.TipoPessoa,
		Nome:               cliente.Nome,
		RelacionamentosCCS: relacionamentos,
		Autorizado:         true,
		Status:             "Sucesso",
	}
	
	// Salvar requisição
	_, err = s.ccsRepo.CriarRequisicaoRelacionamentoCCS(requisicao)
	if err != nil {
		return nil, err
	}
	
	// Buscar a requisição salva com todos os relacionamentos
	reqSalva, err := s.ccsRepo.BuscarRequisicoesRelacionamentoCCS(cpfResponsavel)
	if err != nil {
		return nil, err
	}
	
	// Retornar apenas a requisição recém-criada (a primeira)
	if len(reqSalva) > 0 {
		return []models.RequisicaoRelacionamentoCCS{reqSalva[0]}, nil
	}
	
	return []models.RequisicaoRelacionamentoCCS{{
		CPFCNPJConsulta: cpfCnpj,
		Status:          "Sucesso",
	}}, nil
}

// SolicitarDetalhamento solicita detalhamento de um relacionamento CCS
func (s *CCSService) SolicitarDetalhamento(numeroRequisicao, cpfCnpj, cnpjResponsavel, cnpjParticipante, dataInicioRelacionamento, nomeBancoResponsavel string, idRelacionamento int) ([]map[string]string, error) {
	// Verificar se está no horário permitido
	now := time.Now()
	isWeekday := now.Weekday() >= time.Monday && now.Weekday() <= time.Friday
	
	early := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	late := time.Date(now.Year(), now.Month(), now.Day(), 18, 55, 0, 0, now.Location())
	
	// Se fora do horário, colocar na fila
	if now.Before(early) || now.After(late) || !isWeekday {
		err := s.ccsRepo.AtualizarStatusDetalhamentoCCS(idRelacionamento, "Na fila", false, false, "", "", "", "")
		if err != nil {
			return nil, err
		}
		
		return []map[string]string{
			{
				"banco":  nomeBancoResponsavel,
				"msg":    "Na fila de processamento",
				"status": "pendente",
			},
		}, nil
	}
	
	// Fazer requisição de detalhamento
	url := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/requisitar-detalhamentos?numeros-requisicoes=%s&ids-pessoa=%s&cnpj-responsaveis=%s&cnpj-participantes=%s&datas-inicio=%s",
		numeroRequisicao, cpfCnpj, cnpjResponsavel, cnpjParticipante, dataInicioRelacionamento)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", s.getBasicAuth())
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Ler resposta XML
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Verificar resposta com erro (código 500)
	if resp.StatusCode == 500 {
		// Instituição financeira não responde a detalhamentos
		err = s.ccsRepo.AtualizarStatusDetalhamentoCCS(
			idRelacionamento,
			"IF não detalha",
			false,
			true,
			time.Now().Format(time.RFC3339),
			"",
			"",
			"",
		)
		if err != nil {
			return nil, err
		}
		
		return []map[string]string{
			{
				"banco":  nomeBancoResponsavel,
				"msg":    "Sem detalhamento",
				"status": "falha",
			},
		}, nil
	}
	
	// Converter XML para estrutura
	var requisicaoDetalhamentosXML RequisicaoDetalhamentosXML
	err = xml.Unmarshal(bodyBytes, &requisicaoDetalhamentosXML)
	if err != nil {
		return nil, err
	}
	
	// Obter data de requisição do detalhamento
	var dataRequisicaoDetalhamento string
	if len(requisicaoDetalhamentosXML.RequisicaoDetalhamento) > 0 {
		dataRequisicaoDetalhamento = requisicaoDetalhamentosXML.RequisicaoDetalhamento[0].DataHoraRequisicao
	} else {
		dataRequisicaoDetalhamento = time.Now().Format(time.RFC3339)
	}
	
	// Atualizar status do relacionamento
	err = s.ccsRepo.AtualizarStatusDetalhamentoCCS(
		idRelacionamento,
		"Solicitado. Aguardando...",
		true,
		false,
		dataRequisicaoDetalhamento,
		"",
		"",
		"",
	)
	if err != nil {
		return nil, err
	}
	
	return []map[string]string{
		{
			"banco":  nomeBancoResponsavel,
			"msg":    "Detalhamento Solicitado",
			"status": "sucesso",
		},
	}, nil
}

// ProcessarFilaCCS processa a fila de solicitações de detalhamento CCS
func (s *CCSService) ProcessarFilaCCS() ([]map[string]string, error) {
	// Verificar se está no horário permitido
	now := time.Now()
	isWeekday := now.Weekday() >= time.Monday && now.Weekday() <= time.Friday
	
	early := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	late := time.Date(now.Year(), now.Month(), now.Day(), 18, 55, 0, 0, now.Location())
	
	// Se fora do horário, retornar mensagem
	if now.Before(early) || now.After(late) || !isWeekday {
		return []map[string]string{
			{
				"msg":    "Detalhamento somente pode ser solicitado entre 10h e 19h",
				"status": "falha",
			},
		}, nil
	}
	
	// Buscar relacionamentos na fila
	requisicoes, err := s.ccsRepo.BuscarRelacionamentosNaFila()
	if err != nil {
		return nil, err
	}
	
	resultados := make([]map[string]string, 0)
	
	// Processar cada requisição
	for _, req := range requisicoes {
		for _, relacionamento := range req.RelacionamentosCCS {
			// Fazer requisição de detalhamento
			url := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/requisitar-detalhamentos?numeros-requisicoes=%s&ids-pessoa=%s&cnpj-responsaveis=%s&cnpj-participantes=%s&datas-inicio=%s",
				req.NumeroRequisicao, req.CPFCNPJ, relacionamento.CNPJResponsavel, relacionamento.CNPJParticipante, relacionamento.DataInicioRelacionamento)
			
			httpReq, err := http.NewRequest("GET", url, nil)
			if err != nil {
				continue
			}
			
			httpReq.Header.Add("Accept", "*/*")
			httpReq.Header.Add("Authorization", s.getBasicAuth())
			
			resp, err := s.client.Do(httpReq)
			if err != nil {
				continue
			}
			
			// Ler resposta XML
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			
			if err != nil {
				continue
			}
			
			// Verificar resposta com erro (código 500)
			if resp.StatusCode == 500 {
				// Instituição financeira não responde a detalhamentos
				err = s.ccsRepo.AtualizarStatusDetalhamentoCCS(
					relacionamento.ID,
					"IF não detalha",
					false,
					true,
					time.Now().Format(time.RFC3339),
					"",
					"",
					"",
				)
				if err == nil {
					resultados = append(resultados, map[string]string{
						"banco":  relacionamento.NomeBancoResponsavel,
						"msg":    "Sem detalhamento",
						"status": "falha",
					})
				}
				continue
			}
			
			// Converter XML para estrutura
			var requisicaoDetalhamentosXML RequisicaoDetalhamentosXML
			err = xml.Unmarshal(bodyBytes, &requisicaoDetalhamentosXML)
			if err != nil {
				continue
			}
			
			// Obter data de requisição do detalhamento
			var dataRequisicaoDetalhamento string
			if len(requisicaoDetalhamentosXML.RequisicaoDetalhamento) > 0 {
				dataRequisicaoDetalhamento = requisicaoDetalhamentosXML.RequisicaoDetalhamento[0].DataHoraRequisicao
			} else {
				dataRequisicaoDetalhamento = time.Now().Format(time.RFC3339)
			}
			
			// Atualizar status do relacionamento
			err = s.ccsRepo.AtualizarStatusDetalhamentoCCS(
				relacionamento.ID,
				"Solicitado. Aguardando...",
				true,
				false,
				dataRequisicaoDetalhamento,
				"",
				"",
				"",
			)
			if err == nil {
				resultados = append(resultados, map[string]string{
					"banco":  relacionamento.NomeBancoResponsavel,
					"msg":    "Detalhamento Solicitado",
					"status": "sucesso",
				})
			}
		}
	}
	
	return resultados, nil
}

// ReceberBDVCCS recebe as respostas de detalhamento e BDVs do CCS
func (s *CCSService) ReceberBDVCCS() error {
	// Buscar relacionamentos aguardando resposta
	requisicoes, err := s.ccsRepo.BuscarRelacionamentosAguardandoResposta()
	if err != nil {
		return err
	}
	
	// Processar cada requisição
	for _, req := range requisicoes {
		for _, relacionamento := range req.RelacionamentosCCS {
			// Fazer requisição para obter respostas de detalhamento
			url := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/obter-respostas-detalhamento?numero-requisicao=%s&id-pessoa=%s&cnpj-responsavel=%s&cnpj-participante=%s",
				req.NumeroRequisicao, req.CPFCNPJ, relacionamento.CNPJResponsavel, relacionamento.CNPJParticipante)
			
			httpReq, err := http.NewRequest("GET", url, nil)
			if err != nil {
				continue
			}
			
			httpReq.Header.Add("Accept", "*/*")
			httpReq.Header.Add("Authorization", s.getBasicAuth())
			
			resp, err := s.client.Do(httpReq)
			if err != nil {
				continue
			}
			
			// Ler resposta XML
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			
			if err != nil {
				continue
			}
			
			// Converter XML para estrutura
			var respostaDetalhamentosXML RespostaDetalhamentosXML
			err = xml.Unmarshal(bodyBytes, &respostaDetalhamentosXML)
			if err != nil {
				continue
			}
			
			// Verificar se há respostas
			var resposta bool
			var codigoResposta, codigoIfResposta, nuopResposta string
			
			if len(respostaDetalhamentosXML.RespostaDetalhamento) > 0 {
				resposta = true
				codigoResposta = respostaDetalhamentosXML.RespostaDetalhamento[0].Codigo
				
				if respostaDetalhamentosXML.RespostaDetalhamento[0].CodigoIf != "" {
					codigoIfResposta = respostaDetalhamentosXML.RespostaDetalhamento[0].CodigoIf
				}
				
				if respostaDetalhamentosXML.RespostaDetalhamento[0].Nuop != "" {
					nuopResposta = respostaDetalhamentosXML.RespostaDetalhamento[0].Nuop
				}
				
				// Atualizar status do relacionamento
				err = s.ccsRepo.AtualizarStatusDetalhamentoCCS(
					relacionamento.ID,
					"Concluído",
					true,
					true,
					relacionamento.DataRequisicaoDetalhamento,
					codigoResposta,
					codigoIfResposta,
					nuopResposta,
				)
				if err != nil {
					continue
				}
				
				// Buscar BDVs se houver resposta
				if resposta {
					// Fazer requisição para obter BDVs
					urlBDV := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/obter-bdvs-resposta?numero-controle-resposta=%s", codigoResposta)
					
					httpReqBDV, err := http.NewRequest("GET", urlBDV, nil)
					if err != nil {
						continue
					}
					
					httpReqBDV.Header.Add("Accept", "*/*")
					httpReqBDV.Header.Add("Authorization", s.getBasicAuth())
					
					respBDV, err := s.client.Do(httpReqBDV)
					if err != nil {
						continue
					}
					
					// Ler resposta XML
					bodyBytesBDV, err := ioutil.ReadAll(respBDV.Body)
					respBDV.Body.Close()
					
					if err != nil {
						continue
					}
					
					// Converter XML para estrutura
					var bemDireitoValorsXML BemDireitoValorsXML
					err = xml.Unmarshal(bodyBytesBDV, &bemDireitoValorsXML)
					if err != nil {
						continue
					}
					
					// Processar cada BDV
					for _, bdvXML := range bemDireitoValorsXML.BemDireitoValor {
						// Criar modelo de BDV
						bdv := models.BemDireitoValorCCS{
							CNPJParticipante:  bdvXML.CNPJParticipante,
							Tipo:              bdvXML.Tipo,
							Agencia:           bdvXML.Agencia,
							Conta:             bdvXML.Conta,
							Vinculo:           bdvXML.Vinculo,
							NomePessoa:        bdvXML.NomePessoa,
							DataInicio:        bdvXML.DataInicio,
							DataFim:           bdvXML.DataFim,
							IDRelacionamento:  relacionamento.ID,
						}
						
						// Criar vinculados
						var vinculados []models.VinculadosBDVCCS
						if len(bdvXML.Vinculados.Vinculados) > 0 {
							for _, vincXML := range bdvXML.Vinculados.Vinculados {
								vinculado := models.VinculadosBDVCCS{
									IDPessoa:          vincXML.IDPessoa,
									DataInicio:        vincXML.DataInicio,
									DataFim:           vincXML.DataFim,
									NomePessoa:        vincXML.NomePessoa,
									NomePessoaReceita: vincXML.NomePessoaReceita,
									Tipo:              vincXML.Tipo,
								}
								vinculados = append(vinculados, vinculado)
							}
						}
						
						bdv.Vinculados = vinculados
						
						// Salvar BDV e vinculados
						// Implementação simplificada: na prática, você usaria o repositório
						// para salvar esses dados
					}
				}
			}
		}
	}
	
	return nil
}

// BuscarRequisicoesCCS busca todas as requisições CCS de um CPF responsável
func (s *CCSService) BuscarRequisicoesCCS(cpfResponsavel string) ([]models.RequisicaoRelacionamentoCCS, error) {
	return s.ccsRepo.BuscarRequisicoesRelacionamentoCCS(cpfResponsavel)
}

// ConsultarParticipante consulta informações do participante pelo CNPJ
func (s *CCSService) ConsultarParticipante(cnpj string) (*ParticipanteResponse, error) {
	url := fmt.Sprintf("https://www3.bcb.gov.br/informes/rest/pessoasJuridicas?cnpj=%s", cnpj)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao consultar participante")
	}
	
	var participante ParticipanteResponse
	err = json.NewDecoder(resp.Body).Decode(&participante)
	if err != nil {
		return nil, err
	}
	
	return &participante, nil
}