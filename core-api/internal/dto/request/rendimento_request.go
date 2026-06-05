package request

type RendimentoRequest struct {
	Meses int `json:"meses" binding:"required,gt=0"`
}
