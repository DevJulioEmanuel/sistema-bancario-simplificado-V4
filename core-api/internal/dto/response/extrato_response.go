package response

type ExtratoResponse struct {
	NumeroConta int                 `json:"numeroConta"`
	Titular     string              `json:"titular"`
	SaldoAtual  float64             `json:"saldoAtual"`
	Transacoes  []TransacaoResponse `json:"transacoes"`
}
