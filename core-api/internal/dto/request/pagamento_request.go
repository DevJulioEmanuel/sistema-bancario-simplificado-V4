package request

type PagamentoRequest struct {
	Descricao string  `json:"descricao"`
	Valor     float64 `json:"valor"`
}
