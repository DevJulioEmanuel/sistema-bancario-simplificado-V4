package response

type RendimentoResponse struct {
	NumeroConta    int     `json:"numeroConta"`
	SaldoAtual     float64 `json:"saldoAtual"`
	Meses          int     `json:"meses"`
	TaxaMensal     float64 `json:"taxaMensal"`
	SaldoProjetado float64 `json:"saldoProjetado"`
	Rendimento     float64 `json:"rendimento"`
}
