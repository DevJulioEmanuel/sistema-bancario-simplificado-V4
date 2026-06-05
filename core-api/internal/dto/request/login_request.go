package request

import "shared"

type LoginRequest struct {
	CPF       string           `json:"cpf" binding:"required"`
	Senha     string           `json:"senha" binding:"required"`
	TipoConta shared.TipoConta `json:"tipoConta"`
}
