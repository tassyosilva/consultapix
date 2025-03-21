package models

type RequisicaoRelacionamentoCCS struct {
	ID                  int                   `json:"id" db:"id"`
	DataRequisicao      string                `json:"dataRequisicao" db:"data_requisicao"`
	DataInicioConsulta  string                `json:"dataInicioConsulta" db:"data_inicio_consulta"`
	DataFimConsulta     string                `json:"dataFimConsulta" db:"data_fim_consulta"`
	CPFCNPJConsulta     string                `json:"cpfCnpjConsulta" db:"cpf_cnpj_consulta"`
	NumeroProcesso      string                `json:"numeroProcesso" db:"numero_processo"`
	MotivoBusca         string                `json:"motivoBusca" db:"motivo_busca"`
	CPFResponsavel      string                `json:"cpfResponsavel" db:"cpf_responsavel"`
	Lotacao             string                `json:"lotacao" db:"lotacao"`
	Caso                string                `json:"caso" db:"caso"`
	NumeroRequisicao    string                `json:"numeroRequisicao" db:"numero_requisicao"`
	CPFCNPJ             string                `json:"cpfCnpj" db:"cpf_cnpj"`
	TipoPessoa          string                `json:"tipoPessoa" db:"tipo_pessoa"`
	Nome                string                `json:"nome" db:"nome"`
	RelacionamentosCCS  []RelacionamentoCCS   `json:"relacionamentosCCS,omitempty"`
	Autorizado          bool                  `json:"autorizado" db:"autorizado"`
	CPFAutorizacao      string                `json:"cpfAutorizacao" db:"cpf_autorizacao"`
	NomeAutorizacao     string                `json:"nomeAutorizacao" db:"nome_autorizacao"`
	DataHoraAutorizacao string                `json:"dataHoraAutorizacao" db:"data_hora_autorizacao"`
	TokenAutorizacao    string                `json:"tokenAutorizacao" db:"token_autorizacao"`
	Status              string                `json:"status" db:"status"`
	Detalhamento        bool                  `json:"detalhamento" db:"detalhamento"`
}

type RelacionamentoCCS struct {
	ID                      int                   `json:"id" db:"id"`
	NumeroRequisicao        string                `json:"numeroRequisicao" db:"numero_requisicao"`
	IDPessoa                string                `json:"idPessoa" db:"id_pessoa"`
	NomePessoa              string                `json:"nomePessoa" db:"nome_pessoa"`
	TipoPessoa              string                `json:"tipoPessoa" db:"tipo_pessoa"`
	CNPJResponsavel         string                `json:"cnpjResponsavel" db:"cnpj_responsavel"`
	NumeroBancoResponsavel  string                `json:"numeroBancoResponsavel" db:"numero_banco_responsavel"`
	NomeBancoResponsavel    string                `json:"nomeBancoResponsavel" db:"nome_banco_responsavel"`
	CNPJParticipante        string                `json:"cnpjParticipante" db:"cnpj_participante"`
	NumeroBancoParticipante string                `json:"numeroBancoParticipante" db:"numero_banco_participante"`
	NomeBancoParticipante   string                `json:"nomeBancoParticipante" db:"nome_banco_participante"`
	DataInicioRelacionamento string               `json:"dataInicioRelacionamento" db:"data_inicio_relacionamento"`
	DataFimRelacionamento   string                `json:"dataFimRelacionamento" db:"data_fim_relacionamento"`
	IDRequisicao            int                   `json:"idRequisicao" db:"id_requisicao"`
	DataRequisicaoDetalhamento string             `json:"dataRequisicaoDetalhamento" db:"data_requisicao_detalhamento"`
	StatusDetalhamento      string                `json:"statusDetalhamento" db:"status_detalhamento"`
	RespondeDetalhamento    bool                  `json:"respondeDetalhamento" db:"responde_detalhamento"`
	Resposta                bool                  `json:"resposta" db:"resposta"`
	CodigoResposta          string                `json:"codigoResposta" db:"codigo_resposta"`
	CodigoIfResposta        string                `json:"codigoIfResposta" db:"codigo_if_resposta"`
	NuopResposta            string                `json:"nuopResposta" db:"nuop_resposta"`
	BemDireitoValorCCS      []BemDireitoValorCCS  `json:"bemDireitoValorCCS,omitempty"`
}

type BemDireitoValorCCS struct {
	ID                int                 `json:"id" db:"id"`
	CNPJParticipante  string              `json:"cnpjParticipante" db:"cnpj_participante"`
	Tipo              string              `json:"tipo" db:"tipo"`
	Agencia           string              `json:"agencia" db:"agencia"`
	Conta             string              `json:"conta" db:"conta"`
	Vinculo           string              `json:"vinculo" db:"vinculo"`
	NomePessoa        string              `json:"nomePessoa" db:"nome_pessoa"`
	DataInicio        string              `json:"dataInicio" db:"data_inicio"`
	DataFim           string              `json:"dataFim" db:"data_fim"`
	IDRelacionamento  int                 `json:"idRelacionamento" db:"id_relacionamento"`
	Vinculados        []VinculadosBDVCCS  `json:"vinculados,omitempty"`
}

type VinculadosBDVCCS struct {
	ID                 int    `json:"id" db:"id"`
	IDBDV              int    `json:"idBDV" db:"id_bdv"`
	DataInicio         string `json:"dataInicio" db:"data_inicio"`
	DataFim            string `json:"dataFim" db:"data_fim"`
	IDPessoa           string `json:"idPessoa" db:"id_pessoa"`
	NomePessoa         string `json:"nomePessoa" db:"nome_pessoa"`
	NomePessoaReceita  string `json:"nomePessoaReceita" db:"nome_pessoa_receita"`
	Tipo               string `json:"tipo" db:"tipo"`
}