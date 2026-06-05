package response

type PagamentoResponse struct {
	Mensagem    string  `json:"mensagem"`
	NumeroConta int     `json:"numeroConta"`
	Descricao   string  `json:"descricao"`
	ValorPago   float64 `json:"valorPago"`
	SaldoAtual  float64 `json:"saldoAtual"`
}
