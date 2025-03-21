package models

import "time"

type RequisicaoPix struct {
	ID              int         `json:"id" db:"id"`
	Data            time.Time   `json:"data" db:"data"`
	CPFResponsavel  string      `json:"cpfResponsavel" db:"cpf_responsavel"`
	Lotacao         string      `json:"lotacao" db:"lotacao"`
	Caso            string      `json:"caso" db:"caso"`
	TipoBusca       string      `json:"tipoBusca" db:"tipo_busca"`
	ChaveBusca      string      `json:"chaveBusca" db:"chave_busca"`
	MotivoBusca     string      `json:"motivoBusca" db:"motivo_busca"`
	Resultado       string      `json:"resultado" db:"resultado"`
	Vinculos        interface{} `json:"vinculos" db:"vinculos"`
	Chaves          []ChavePix  `json:"chaves,omitempty"`
	Autorizado      bool        `json:"autorizado" db:"autorizado"`
	CPFAutorizacao  string      `json:"cpfAutorizacao" db:"cpf_autorizacao"`
	NomeAutorizacao string      `json:"nomeAutorizacao" db:"nome_autorizacao"`
	DataHoraAutorizacao string  `json:"dataHoraAutorizacao" db:"data_hora_autorizacao"`
	TokenAutorizacao string     `json:"tokenAutorizacao" db:"token_autorizacao"`
}

type ChavePix struct {
	ID                     int               `json:"id" db:"id"`
	Chave                  string            `json:"chave" db:"chave"`
	TipoChave              string            `json:"tipoChave" db:"tipo_chave"`
	Status                 string            `json:"status" db:"status"`
	DataAberturaReivindicacao string         `json:"dataAberturaReivindicacao" db:"data_abertura_reivindicacao"`
	CPFCNPJ                string            `json:"cpfCnpj" db:"cpf_cnpj"`
	NomeProprietario       string            `json:"nomeProprietario" db:"nome_proprietario"`
	NomeFantasia           string            `json:"nomeFantasia" db:"nome_fantasia"`
	Participante           string            `json:"participante" db:"participante"`
	Agencia                string            `json:"agencia" db:"agencia"`
	NumeroConta            string            `json:"numeroConta" db:"numero_conta"`
	TipoConta              string            `json:"tipoConta" db:"tipo_conta"`
	DataAberturaConta      string            `json:"dataAberturaConta" db:"data_abertura_conta"`
	ProprietarioDaChaveDesde string          `json:"proprietarioDaChaveDesde" db:"proprietario_da_chave_desde"`
	DataCriacao            string            `json:"dataCriacao" db:"data_criacao"`
	UltimaModificacao      string            `json:"ultimaModificacao" db:"ultima_modificacao"`
	NumeroBanco            string            `json:"numeroBanco" db:"numero_banco"`
	NomeBanco              string            `json:"nomeBanco" db:"nome_banco"`
	CPFCNPJBusca           string            `json:"cpfCnpjBusca" db:"cpf_cnpj_busca"`
	NomeProprietarioBusca  string            `json:"nomeProprietarioBusca" db:"nome_proprietario_busca"`
	EventosVinculo         []EventoChavePix  `json:"eventosVinculo,omitempty"`
	IDRequisicao           int               `json:"idRequisicao" db:"id_requisicao"`
}

type EventoChavePix struct {
	ID                int     `json:"id" db:"id"`
	TipoEvento        string  `json:"tipoEvento" db:"tipo_evento"`
	MotivoEvento      string  `json:"motivoEvento" db:"motivo_evento"`
	DataEvento        string  `json:"dataEvento" db:"data_evento"`
	Chave             string  `json:"chave" db:"chave"`
	TipoChave         string  `json:"tipoChave" db:"tipo_chave"`
	CPFCNPJ           string  `json:"cpfCnpj" db:"cpf_cnpj"`
	NomeProprietario  string  `json:"nomeProprietario" db:"nome_proprietario"`
	NomeFantasia      string  `json:"nomeFantasia" db:"nome_fantasia"`
	Participante      string  `json:"participante" db:"participante"`
	Agencia           string  `json:"agencia" db:"agencia"`
	NumeroConta       string  `json:"numeroConta" db:"numero_conta"`
	TipoConta         string  `json:"tipoConta" db:"tipo_conta"`
	DataAberturaConta string  `json:"dataAberturaConta" db:"data_abertura_conta"`
	NumeroBanco       string  `json:"numeroBanco" db:"numero_banco"`
	NomeBanco         string  `json:"nomeBanco" db:"nome_banco"`
	IDChave           int     `json:"idChave" db:"id_chave"`
}