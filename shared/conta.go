package shared

import (
	"sync"
)

type TipoConta int

const (
	TipoCorrente TipoConta = 1
	TipoPoupanca TipoConta = 2
	TaxaMensal             = 0.005
)

type Conta struct {
	Numero     int          `json:"numero" gorm:"primaryKey"`
	Saldo      float64      `json:"saldo" gorm:"not null;default:0"`
	TitularCPF string       `json:"-" gorm:"not null"`
	Titular    *Cliente     `json:"titular" gorm:"foreignKey:TitularCPF"`
	Tipo       TipoConta    `json:"tipo" gorm:"not null"`
	Historico  []Transacao  `json:"historico" gorm:"foreignKey:ContaNumero"`
	Mu         sync.RWMutex `gorm:"-"`
}
