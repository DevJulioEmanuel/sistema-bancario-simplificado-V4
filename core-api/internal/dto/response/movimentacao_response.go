package response

type MovimentacaoResponse struct {
	Mensagem         string  `json:"mensagem"`
	NumeroConta      int     `json:"numeroConta"`
	TipoOperacao     string  `json:"tipoOperacao"`
	ValorMovimentado float64 `json:"valorMovimentado"`
	SaldoAtual       float64 `json:"saldoAtual"`
}
