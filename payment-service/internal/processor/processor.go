package processor

import (
	"errors"
	"fmt"
	"log"
	"shared"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"payment-service/internal/publisher"

	"gorm.io/gorm"
)

type Processor struct {
	db *gorm.DB
}

func NewProcessor(db *gorm.DB) *Processor {
	return &Processor{db: db}
}

func (p *Processor) ProcessarDeposito(ch *amqp.Channel, transacao shared.TransacaoEvent) error {

	result := p.db.Model(&shared.Conta{}).
		Where("numero = ?", transacao.NumeroConta).
		Update("saldo", gorm.Expr("saldo + ?", transacao.Valor))

	if result.Error != nil {
		return result.Error
	}

	novaTransacao := shared.Transacao{
		ContaNumero: transacao.NumeroConta,
		Tipo:        "deposito",
		Descricao:   "Deposito",
		Valor:       +transacao.Valor,
		Data:        time.Now(),
	}

	if err := p.db.Create(&novaTransacao).Error; err != nil {
		log.Println("⚠️ Erro ao gravar histórico de depósito no banco:", err)
	}

	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)

	publisher.EnviarRespostaDeSaldo(ch, transacao.NumeroConta, conta.Saldo)
	publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Sucesso: deposito de R$ "+fmt.Sprintf("%.2f", transacao.Valor)+" realizado com sucesso.", true)
	return nil

}

func (p *Processor) ProcessarSaque(ch *amqp.Channel, transacao shared.TransacaoEvent) error {
	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	if conta.Saldo < transacao.Valor {
		publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Falha: saldo insuficiente para saque.", false)
		return errors.New("saldo insuficiente")
	}
	result := p.db.Model(&shared.Conta{}).
		Where("numero = ?", transacao.NumeroConta).
		Update("saldo", gorm.Expr("saldo - ?", transacao.Valor))

	if result.Error != nil {
		return result.Error
	}

	novaTransacao := shared.Transacao{
		ContaNumero: transacao.NumeroConta,
		Tipo:        "saque",
		Descricao:   "Saque",
		Valor:       -transacao.Valor,
		Data:        time.Now(),
	}

	if err := p.db.Create(&novaTransacao).Error; err != nil {
		log.Println("⚠️ Erro ao gravar histórico de saque no banco:", err)
	}

	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	publisher.EnviarRespostaDeSaldo(ch, transacao.NumeroConta, conta.Saldo)
	publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Sucesso: saque de R$ "+fmt.Sprintf("%.2f", transacao.Valor)+" realizado com sucesso.", true)
	return nil
}

func (p *Processor) ProcessarPagamento(ch *amqp.Channel, transacao shared.TransacaoEvent) error {
	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	if conta.Saldo < transacao.Valor {
		publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Falha: saldo insuficiente para pagamento.", false)
		return errors.New("saldo insuficiente")
	}
	result := p.db.Model(&shared.Conta{}).
		Where("numero = ?", transacao.NumeroConta).
		Update("saldo", gorm.Expr("saldo - ?", transacao.Valor))

	if result.Error != nil {
		return result.Error
	}

	novaTransacao := shared.Transacao{
		ContaNumero: transacao.NumeroConta,
		Tipo:        "PAGAMENTO",
		Descricao:   transacao.Descricao,
		Valor:       -transacao.Valor,
		Data:        time.Now(),
	}

	if err := p.db.Create(&novaTransacao).Error; err != nil {
		log.Println("⚠️ Erro ao gravar histórico de pagamento no banco:", err)
	}

	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	publisher.EnviarRespostaDeSaldo(ch, transacao.NumeroConta, conta.Saldo)
	publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Sucesso: pagamento de R$ "+fmt.Sprintf("%.2f", transacao.Valor)+" realizado com sucesso.", true)
	return nil
}

func (p *Processor) ProcessarTransferencia(ch *amqp.Channel, transacao shared.TransacaoEvent) error {
	var conta shared.Conta
	p.db.Where("numero = ?", transacao.NumeroConta).First(&conta)
	if conta.Saldo < transacao.Valor {
		publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Falha: saldo insuficiente para transferência.", false)
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

	if err != nil {
		return err
	}

	historicoOrigem := shared.Transacao{
		ContaNumero: transacao.NumeroConta,
		Tipo:        "TRANSFERENCIA",
		Descricao:   fmt.Sprintf("Transf. enviada para Conta %d", transacao.NumContaDestino),
		Valor:       -transacao.Valor,
		Data:        time.Now(),
	}
	if err := p.db.Create(&historicoOrigem).Error; err != nil {
		log.Println("⚠️ Erro ao gravar histórico de transferência na origem:", err)
	}

	historicoDestino := shared.Transacao{
		ContaNumero: transacao.NumContaDestino,
		Tipo:        "TRANSFERENCIA",
		Descricao:   fmt.Sprintf("Transf. recebida da Conta %d", transacao.NumeroConta),
		Valor:       transacao.Valor,
		Data:        time.Now(),
	}
	if err := p.db.Create(&historicoDestino).Error; err != nil {
		log.Println("⚠️ Erro ao gravar histórico de transferência no destino:", err)
	}

	var contaOrigem shared.Conta
	var contaDestino shared.Conta

	p.db.Where("numero = ?", transacao.NumeroConta).First(&contaOrigem)
	p.db.Where("numero = ?", transacao.NumContaDestino).First(&contaDestino)

	publisher.EnviarRespostaDeSaldo(ch, transacao.NumeroConta, contaOrigem.Saldo)

	publisher.EnviarRespostaDeSaldo(ch, transacao.NumContaDestino, contaDestino.Saldo)
	publisher.EnviarNotificacao(ch, transacao.NumeroConta, "Sucesso: transferência de R$ "+fmt.Sprintf("%.2f", transacao.Valor)+" realizada com sucesso.", true)
	publisher.EnviarNotificacao(ch, transacao.NumContaDestino, "Sucesso: transferência de R$ "+fmt.Sprintf("%.2f", transacao.Valor)+" recebida com sucesso.", true)

	return nil
}
