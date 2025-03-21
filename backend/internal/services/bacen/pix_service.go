package bacen

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tassyosilva/consultapix/internal/config"
	"github.com/tassyosilva/consultapix/internal/database/models"
	"github.com/tassyosilva/consultapix/internal/repository"
)

type PixService struct {
	config   *config.Config
	pixRepo  *repository.PixRepository
	client   *http.Client
}

type ParticipanteResponse struct {
	CodigoCompensacao int    `json:"codigoCompensacao"`
	Nome              string `json:"nome"`
}

type ChavePixResponse struct {
	Chave                     string              `json:"chave"`
	TipoChave                 string              `json:"tipoChave"`
	Status                    string              `json:"status"`
	DataAberturaReivindicacao string              `json:"dataAberturaReivindicacao"`
	CPFCNPJ                   string              `json:"cpfCnpj"`
	NomeProprietario          string              `json:"nomeProprietario"`
	NomeFantasia              string              `json:"nomeFantasia"`
	Participante              string              `json:"participante"`
	Agencia                   string              `json:"agencia"`
	NumeroConta               string              `json:"numeroConta"`
	TipoConta                 string              `json:"tipoConta"`
	DataAberturaConta         string              `json:"dataAberturaConta"`
	ProprietarioDaChaveDesde  string              `json:"proprietarioDaChaveDesde"`
	DataCriacao               string              `json:"dataCriacao"`
	UltimaModificacao         string              `json:"ultimaModificacao"`
	EventosVinculo            []EventoVinculoResponse `json:"eventosVinculo"`
	// Campos adicionados na aplicação
	NumeroBanco              string               `json:"numerobanco"`
	NomeBanco                string               `json:"nomebanco"`
	CPFCNPJBusca             string               `json:"cpfCnpjBusca"`
	NomeProprietarioBusca    string               `json:"nomeProprietarioBusca"`
}

type EventoVinculoResponse struct {
	TipoEvento        string `json:"tipoEvento"`
	MotivoEvento      string `json:"motivoEvento"`
	DataEvento        string `json:"dataEvento"`
	Chave             string `json:"chave"`
	TipoChave         string `json:"tipoChave"`
	CPFCNPJ           string `json:"cpfCnpj"`
	NomeProprietario  string `json:"nomeProprietario"`
	NomeFantasia      string `json:"nomeFantasia"`
	Participante      string `json:"participante"`
	Agencia           string `json:"agencia"`
	NumeroConta       string `json:"numeroConta"`
	TipoConta         string `json:"tipoConta"`
	DataAberturaConta string `json:"dataAberturaConta"`
	// Campos adicionados na aplicação
	NumeroBanco       string `json:"numerobanco"`
	NomeBanco         string `json:"nomebanco"`
}

type VinculosPixResponse struct {
	VinculosPix []ChavePixResponse `json:"vinculosPix"`
}

func NewPixService(cfg *config.Config) *PixService {
	return &PixService{
		config:  cfg,
		pixRepo: repository.NewPixRepository(),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// getBasicAuth retorna o cabeçalho de autenticação básica
func (s *PixService) getBasicAuth() string {
	credentials := fmt.Sprintf("%s:%s", s.config.BacenUsername, s.config.BacenPassword)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}

// ConsultarParticipante consulta informações do participante pelo CNPJ
func (s *PixService) ConsultarParticipante(cnpj string) (*ParticipanteResponse, error) {
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
		return nil, errors.New("falha ao consultar participante")
	}
	
	var participante ParticipanteResponse
	err = json.NewDecoder(resp.Body).Decode(&participante)
	if err != nil {
		return nil, err
	}
	
	return &participante, nil
}

// ConsultarChavePix consulta informações de uma chave PIX
func (s *PixService) ConsultarChavePix(chave, motivo string, cpfResponsavel, lotacao, caso string) ([]ChavePixResponse, error) {
	url := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/consultar-vinculo-pix?chave=%s&motivo=%s", chave, motivo)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", s.getBasicAuth())
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao consultar chave PIX: %d", resp.StatusCode)
	}
	
	var chaveResp ChavePixResponse
	err = json.NewDecoder(resp.Body).Decode(&chaveResp)
	if err != nil {
		return nil, err
	}
	
	// Verificar se a chave existe
	if chaveResp.Chave == "" {
		// Registrar requisição sem chave encontrada
		req := &models.RequisicaoPix{
			Data:           time.Now(),
			CPFResponsavel: cpfResponsavel,
			Lotacao:        lotacao,
			Caso:           caso,
			TipoBusca:      "chave",
			ChaveBusca:     chave,
			MotivoBusca:    motivo,
			Resultado:      "Chave não encontrada",
			Autorizado:     true,
		}
		
		_, err = s.pixRepo.CriarRequisicaoPix(req)
		if err != nil {
			return nil, err
		}
		
		return []ChavePixResponse{}, nil
	}
	
	// Buscar informações do banco participante
	participante, err := s.ConsultarParticipante(chaveResp.Participante)
	if err != nil {
		chaveResp.NumeroBanco = "000"
		chaveResp.NomeBanco = "BANCO NÃO INFORMADO"
	} else {
		// Formatar número do banco com 3 dígitos
		chaveResp.NumeroBanco = fmt.Sprintf("%03d", participante.CodigoCompensacao)
		chaveResp.NomeBanco = participante.Nome
	}
	
	// Processar eventos de vínculo
	for i := range chaveResp.EventosVinculo {
		participante, err := s.ConsultarParticipante(chaveResp.EventosVinculo[i].Participante)
		if err != nil {
			chaveResp.EventosVinculo[i].NumeroBanco = "000"
			chaveResp.EventosVinculo[i].NomeBanco = "BANCO NÃO INFORMADO"
		} else {
			// Formatar número do banco com 3 dígitos
			chaveResp.EventosVinculo[i].NumeroBanco = fmt.Sprintf("%03d", participante.CodigoCompensacao)
			chaveResp.EventosVinculo[i].NomeBanco = participante.Nome
		}
	}
	
	// Se o status estiver vazio, definir como INATIVO
	if chaveResp.Status == "" {
		chaveResp.Status = "INATIVO"
	}
	
	// Se o CPF/CNPJ estiver vazio, tentar obter dos eventos
	if chaveResp.CPFCNPJ == "" {
		if len(chaveResp.EventosVinculo) > 0 && chaveResp.EventosVinculo[0].NomeProprietario != "" {
			chaveResp.NomeProprietario = chaveResp.EventosVinculo[0].NomeProprietario
			chaveResp.CPFCNPJ = chaveResp.EventosVinculo[0].CPFCNPJ
		} else {
			chaveResp.NomeProprietario = "NOME NÃO INFORMADO"
			chaveResp.CPFCNPJ = "CPF/CNPJ NÃO INFORMADO"
		}
	}
	
	// Salvar a requisição no banco de dados
	requisicaoPix := &models.RequisicaoPix{
		Data:           time.Now(),
		CPFResponsavel: cpfResponsavel,
		Lotacao:        lotacao,
		Caso:           caso,
		TipoBusca:      "chave",
		ChaveBusca:     chave,
		MotivoBusca:    motivo,
		Resultado:      "Sucesso",
		Vinculos:       chaveResp,
		Autorizado:     true,
		Chaves: []models.ChavePix{
			{
				Chave:                     chaveResp.Chave,
				TipoChave:                 chaveResp.TipoChave,
				Status:                    chaveResp.Status,
				DataAberturaReivindicacao: chaveResp.DataAberturaReivindicacao,
				CPFCNPJ:                   chaveResp.CPFCNPJ,
				NomeProprietario:          chaveResp.NomeProprietario,
				NomeFantasia:              chaveResp.NomeFantasia,
				Participante:              chaveResp.Participante,
				Agencia:                   chaveResp.Agencia,
				NumeroConta:               chaveResp.NumeroConta,
				TipoConta:                 chaveResp.TipoConta,
				DataAberturaConta:         chaveResp.DataAberturaConta,
				ProprietarioDaChaveDesde:  chaveResp.ProprietarioDaChaveDesde,
				DataCriacao:               chaveResp.DataCriacao,
				UltimaModificacao:         chaveResp.UltimaModificacao,
				NumeroBanco:               chaveResp.NumeroBanco,
				NomeBanco:                 chaveResp.NomeBanco,
				CPFCNPJBusca:              chaveResp.CPFCNPJBusca,
				NomeProprietarioBusca:     chaveResp.NomeProprietarioBusca,
				EventosVinculo:            convertEventosVinculo(chaveResp.EventosVinculo),
			},
		},
	}
	
	_, err = s.pixRepo.CriarRequisicaoPix(requisicaoPix)
	if err != nil {
		return nil, err
	}
	
	return []ChavePixResponse{chaveResp}, nil
}

// ConsultarPorCPFCNPJ consulta todas as chaves PIX associadas a um CPF/CNPJ
func (s *PixService) ConsultarPorCPFCNPJ(cpfCnpj, motivo, cpfResponsavel, lotacao, caso string) (interface{}, error) {
	url := fmt.Sprintf("https://www3.bcb.gov.br/bc_ccs/rest/consultar-vinculos-pix?cpfCnpj=%s&motivo=%s", cpfCnpj, motivo)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", s.getBasicAuth())
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// Armazenar falha na requisição
		errReq := &models.RequisicaoPix{
			Data:           time.Now(),
			CPFResponsavel: cpfResponsavel,
			Lotacao:        lotacao,
			Caso:           caso,
			TipoBusca:      "cpf/cnpj",
			ChaveBusca:     cpfCnpj,
			MotivoBusca:    motivo,
			Autorizado:     true,
			Resultado:      "Erro no processamento da Solicitação",
		}
		
		if resp.StatusCode == http.StatusBadRequest {
			var errorResp struct {
				Message string `json:"message"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Message == "0002 - ERRO_CPF_CNPJ_INVALIDO" {
				errReq.Resultado = "CPF/CNPJ não encontrado"
			}
		}
		
		_, _ = s.pixRepo.CriarRequisicaoPix(errReq)
		return []string{errReq.Resultado}, nil
	}
	
	var vinculosResp VinculosPixResponse
	err = json.NewDecoder(resp.Body).Decode(&vinculosResp)
	if err != nil {
		return nil, err
	}
	
	// Verificar se existem vínculos
	if len(vinculosResp.VinculosPix) == 0 {
		req := &models.RequisicaoPix{
			Data:           time.Now(),
			CPFResponsavel: cpfResponsavel,
			Lotacao:        lotacao,
			Caso:           caso,
			TipoBusca:      "cpf/cnpj",
			ChaveBusca:     cpfCnpj,
			MotivoBusca:    motivo,
			Autorizado:     true,
			Resultado:      "Nenhuma Chave PIX encontrada",
		}
		
		_, err = s.pixRepo.CriarRequisicaoPix(req)
		if err != nil {
			return nil, err
		}
		
		return []string{"Nenhuma Chave PIX encontrada"}, nil
	}
	
	// Ordenar as chaves para que as primeiras contenham CPF/CNPJ e Nome
	// (Implementação simplificada da ordenação)
	
	// Processar cada chave para adicionar informações de banco
	var nomeProprietarioBusca string
	if len(vinculosResp.VinculosPix) > 0 && vinculosResp.VinculosPix[0].NomeProprietario != "" {
		nomeProprietarioBusca = vinculosResp.VinculosPix[0].NomeProprietario
	} else {
		nomeProprietarioBusca = "NOME NÃO INFORMADO"
	}
	
	// Processar cada chave
	chavesModels := make([]models.ChavePix, len(vinculosResp.VinculosPix))
	for i, chave := range vinculosResp.VinculosPix {
		// Buscar informações do banco participante
		participante, err := s.ConsultarParticipante(chave.Participante)
		if err != nil {
			chave.NumeroBanco = "000"
			chave.NomeBanco = "BANCO NÃO INFORMADO"
		} else {
			chave.NumeroBanco = fmt.Sprintf("%03d", participante.CodigoCompensacao)
			chave.NomeBanco = participante.Nome
		}
		
		// Adicionar informações de busca
		chave.CPFCNPJBusca = cpfCnpj
		chave.NomeProprietarioBusca = nomeProprietarioBusca
		
		// Processar eventos
		for j := range chave.EventosVinculo {
			participante, err := s.ConsultarParticipante(chave.EventosVinculo[j].Participante)
			if err != nil {
				chave.EventosVinculo[j].NumeroBanco = "000"
				chave.EventosVinculo[j].NomeBanco = "BANCO NÃO INFORMADO"
			} else {
				chave.EventosVinculo[j].NumeroBanco = fmt.Sprintf("%03d", participante.CodigoCompensacao)
				chave.EventosVinculo[j].NomeBanco = participante.Nome
			}
		}
		
		// Verificar status
		if chave.Status == "" {
			chave.Status = "INATIVO"
		}
		
		// Converter para modelo interno
		chavesModels[i] = models.ChavePix{
			Chave:                     chave.Chave,
			TipoChave:                 chave.TipoChave,
			Status:                    chave.Status,
			DataAberturaReivindicacao: chave.DataAberturaReivindicacao,
			CPFCNPJ:                   chave.CPFCNPJ,
			NomeProprietario:          chave.NomeProprietario,
			NomeFantasia:              chave.NomeFantasia,
			Participante:              chave.Participante,
			Agencia:                   chave.Agencia,
			NumeroConta:               chave.NumeroConta,
			TipoConta:                 chave.TipoConta,
			DataAberturaConta:         chave.DataAberturaConta,
			ProprietarioDaChaveDesde:  chave.ProprietarioDaChaveDesde,
			DataCriacao:               chave.DataCriacao,
			UltimaModificacao:         chave.UltimaModificacao,
			NumeroBanco:               chave.NumeroBanco,
			NomeBanco:                 chave.NomeBanco,
			CPFCNPJBusca:              chave.CPFCNPJBusca,
			NomeProprietarioBusca:     chave.NomeProprietarioBusca,
			EventosVinculo:            convertEventosVinculo(chave.EventosVinculo),
		}
	}
	
	// Salvar requisição
	req := &models.RequisicaoPix{
		Data:           time.Now(),
		CPFResponsavel: cpfResponsavel,
		Lotacao:        lotacao,
		Caso:           caso,
		TipoBusca:      "cpf/cnpj",
		ChaveBusca:     cpfCnpj,
		MotivoBusca:    motivo,
		Resultado:      "Sucesso",
		Vinculos:       vinculosResp.VinculosPix,
		Autorizado:     true,
		Chaves:         chavesModels,
	}
	
	_, err = s.pixRepo.CriarRequisicaoPix(req)
	if err != nil {
		return nil, err
	}
	
	return vinculosResp.VinculosPix, nil
}

// Converter eventos de vínculo para modelo interno
func convertEventosVinculo(eventos []EventoVinculoResponse) []models.EventoChavePix {
	result := make([]models.EventoChavePix, len(eventos))
	for i, evento := range eventos {
		result[i] = models.EventoChavePix{
			TipoEvento:        evento.TipoEvento,
			MotivoEvento:      evento.MotivoEvento,
			DataEvento:        evento.DataEvento,
			Chave:             evento.Chave,
			TipoChave:         evento.TipoChave,
			CPFCNPJ:           evento.CPFCNPJ,
			NomeProprietario:  evento.NomeProprietario,
			NomeFantasia:      evento.NomeFantasia,
			Participante:      evento.Participante,
			Agencia:           evento.Agencia,
			NumeroConta:       evento.NumeroConta,
			TipoConta:         evento.TipoConta,
			DataAberturaConta: evento.DataAberturaConta,
			NumeroBanco:       evento.NumeroBanco,
			NomeBanco:         evento.NomeBanco,
		}
	}
	return result
}