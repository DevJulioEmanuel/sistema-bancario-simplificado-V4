package response

type ContaResponse struct {
	Numero  int     `json:"numero"`
	Titular string  `json:"titular"`
	Saldo   float64 `json:"saldo"`
	Tipo    int     `json:"tipo"`
	Limite  float64 `json:"limite"`
	Taxa    float64 `json:"taxa"`
}
