package model

import "time"

type Notificacao struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ContaNum  int       `gorm:"index" json:"conta_num"`
	Mensagem  string    `json:"mensagem"`
	Sucesso   bool      `json:"sucesso"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"data"`
}
