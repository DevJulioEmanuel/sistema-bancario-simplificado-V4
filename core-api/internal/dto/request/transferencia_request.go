package request

type TransferenciaRequest struct {
	NumDestino int     `json:"numDestino" binding:"required"`
	Valor      float64 `json:"valor" binding:"required,gt=0"`
}
