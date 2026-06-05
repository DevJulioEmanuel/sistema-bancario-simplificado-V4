package request

type ValorRequest struct {
	Valor float64 `json:"valor" binding:"required,gt=0"`
}
