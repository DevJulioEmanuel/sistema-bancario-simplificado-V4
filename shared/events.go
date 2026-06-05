package shared

const (
	FilaNotificacao   = "notificacao"
	FilaDeposito      = "deposito"
	FilaSaque         = "saque"
	FilaTransferencia = "transferencia"
	FilaPagamento     = "pagamento"
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
