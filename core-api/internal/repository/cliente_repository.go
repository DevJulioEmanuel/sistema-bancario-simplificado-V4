package repository

import (
	"shared"

	"gorm.io/gorm"
)

type ClienteRepository struct {
	db *gorm.DB
}

func NewClienteRepository(
	db *gorm.DB,
) *ClienteRepository {
	return &ClienteRepository{
		db: db,
	}
}

func (r *ClienteRepository) Salvar(cliente *shared.Cliente) error {
	result := r.db.Create(cliente)
	return result.Error
}

func (r *ClienteRepository) BuscarPorCPF(cpf string) (*shared.Cliente, bool) {
	var cliente shared.Cliente
	result := r.db.Where("cpf = ?", cpf).First(&cliente)
	if result.Error != nil {
		return nil, false
	}
	return &cliente, true
}
