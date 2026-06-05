package shared

const (
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
}
