package shared

import (
	"fmt"
	"time"
)

type Transacao struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ContaNumero int       `json:"-" gorm:"not null"`
	Tipo        string    `json:"tipo" gorm:"not null"`
	Descricao   string    `json:"descricao"`
	Valor       float64   `json:"valor" gorm:"not null"`
	Data        time.Time `json:"data"`
}

func (t Transacao) String() string {
	sinal := "+"

	if t.Valor < 0 {
		sinal = ""
	}

	return fmt.Sprintf(
		"[%s] %s: %sR$ %.2f",
		t.Data.Format("02/01/2006 15:04"),
		t.Descricao,
		sinal,
		t.Valor,
	)
}
