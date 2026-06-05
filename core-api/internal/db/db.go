package db

import (
	"banco-api/internal/model"
	"shared"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDB(caminho string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(caminho), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&shared.Cliente{}, &shared.Conta{}, &shared.Transacao{}, &model.Notificacao{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
