package shared

type Cliente struct {
	Nome  string `json:"nome" gorm:"not null"`
	CPF   string `json:"cpf" gorm:"primaryKey"`
	Senha string `json:"-" gorm:"not null"`
}
