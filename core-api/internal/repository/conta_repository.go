package repository

import (
	"shared"

	"gorm.io/gorm"
)

type ContaRepository struct {
	db *gorm.DB
}

func NewContaRepository(
	db *gorm.DB,
) *ContaRepository {
	return &ContaRepository{
		db: db,
	}
}

func (r *ContaRepository) Salvar(conta *shared.Conta) error {
	result := r.db.Create(conta)
	return result.Error
}

func (r *ContaRepository) BuscarPorNumero(numero int) (*shared.Conta, bool) {
	var conta shared.Conta
	result := r.db.Preload("Titular").Where("numero = ?", numero).First(&conta)
	if result.Error != nil {
		return nil, false
	}
	return &conta, true
}

func (r *ContaRepository) GerarNumeroConta(tipo shared.TipoConta) int {
	var ultima shared.Conta
	r.db.Where("tipo = ?", tipo).Order("numero desc").First(&ultima)
	if ultima.Numero == 0 {
		if tipo == shared.TipoCorrente {
			return 1001
		}
		return 5001
	}
	return ultima.Numero + 1
}

func (r *ContaRepository) BuscarPorClienteETipo(
	cpf string,
	tipo shared.TipoConta,
) (*shared.Conta, bool) {

	var conta shared.Conta
	result := r.db.Preload("Titular").Where("titular_cpf = ? AND tipo = ?", cpf, tipo).First(&conta)
	if result.Error != nil {
		return nil, false
	}
	return &conta, true
}
