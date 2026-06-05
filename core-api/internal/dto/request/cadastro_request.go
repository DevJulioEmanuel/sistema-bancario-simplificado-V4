package request

type CadastroRequest struct {
	Nome  string `json:"nome" binding:"required"`
	CPF   string `json:"cpf" binding:"required"`
	Senha string `json:"senha" binding:"required"`
	Tipo  int    `json:"tipoConta" binding:"required,oneof=1 2"`
}
