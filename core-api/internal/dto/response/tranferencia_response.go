package response

type TransferenciaResponse struct {
	Mensagem         string  `json:"mensagem"`
	ContaOrigem      int     `json:"contaOrigem"`
	ContaDestino     int     `json:"contaDestino"`
	ValorTransferido float64 `json:"valorTransferido"`
	SaldoAtualOrigem float64 `json:"saldoAtualOrigem"`
}
