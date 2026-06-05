package processor

import (
	"errors"
	"shared"

	"gorm.io/gorm"
)

type Processor struct {
	db *gorm.DB
}

func NewProcessor(db *gorm.DB) *Processor {
	return &Processor{db: db}
}

func (p *Processor) ProcessarDeposito(transacao shared.TransacaoEvent) error {

	result := p.db.Model(&shared.Conta{}).
		Where("numero = ?", transacao.NumeroConta).
		Update("saldo", gorm.Expr("saldo + ?", transacao.Valor))
	return result.Error
}

func (p *Processor) ProcessarSaque(transacao shared.TransacaoEvent) error {
	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	if conta.Saldo < transacao.Valor {
		return errors.New("saldo insuficiente")
	}
	result := p.db.Model(&shared.Conta{}).
		Where("numero = ?", transacao.NumeroConta).
		Update("saldo", gorm.Expr("saldo - ?", transacao.Valor))
	return result.Error
}

func (p *Processor) ProcessarPagamento(transacao shared.TransacaoEvent) error {
	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	if conta.Saldo < transacao.Valor {
		return errors.New("saldo insuficiente")
	}
	result := p.db.Model(&shared.Conta{}).
		Where("numero = ?", transacao.NumeroConta).
		Update("saldo", gorm.Expr("saldo - ?", transacao.Valor))
	return result.Error
}

func (p *Processor) ProcessarTransferencia(transacao shared.TransacaoEvent) error {
	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	if conta.Saldo < transacao.Valor {
		return errors.New("saldo insuficiente")
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&shared.Conta{}).
			Where("numero = ?", transacao.NumeroConta).
			Update("saldo", gorm.Expr("saldo - ?", transacao.Valor)).Error; err != nil {
			return err
		}
		if err := tx.Model(&shared.Conta{}).
			Where("numero = ?", transacao.NumContaDestino).
			Update("saldo", gorm.Expr("saldo + ?", transacao.Valor)).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
