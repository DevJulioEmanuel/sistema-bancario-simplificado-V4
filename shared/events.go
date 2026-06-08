package shared

const (
	FilaNotificacao = "notificacao"
	FilaTransacoes  = "transacoes_conta"
)

type TransacaoEvent struct {
	Tipo            string  `json:"tipoOperacao"`
	NumeroConta     int     `json:"numero"`
	Valor           float64 `json:"valor"`
	NumContaDestino int     `json:"numDestino"`
	Descricao       string  `json:"descricao"`
}

type NotificacaoEvent struct {
	ContaNum int    `json:"conta_num"`
	Mensagem string `json:"mensagem"`
	Sucesso  bool   `json:"sucesso"`
}
