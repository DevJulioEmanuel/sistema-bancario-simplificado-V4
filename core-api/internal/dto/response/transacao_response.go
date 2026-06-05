package response

type TransacaoResponse struct {
	Tipo      string  `json:"tipo"`
	Valor     float64 `json:"valor"`
	Descricao string  `json:"descricao"`
	Data      string  `json:"data"`
}
