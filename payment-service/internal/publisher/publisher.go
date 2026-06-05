package publisher

import (
	"encoding/json"
	"log"
	"shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RespostaFaturamento struct {
	NumeroConta int     `json:"numero_conta"`
	NovoSaldo   float64 `json:"novo_saldo"`
}

func EnviarRespostaDeSaldo(ch *amqp.Channel, numeroConta int, novoSaldo float64) {
	q, err := ch.QueueDeclare(
		"fila_resposta_pagamentos",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Erro ao declarar fila de resposta no payment-service: %v", err)
		return
	}

	dados := RespostaFaturamento{
		NumeroConta: numeroConta,
		NovoSaldo:   novoSaldo,
	}
	body, err := json.Marshal(dados)
	if err != nil {
		log.Printf("Erro ao converter struct para JSON: %v", err)
		return
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Erro ao publicar resposta no RabbitMQ: %v", err)
	}
}

func EnviarNotificacao(canal *amqp.Channel, contaNum int, mensagem string, sucesso bool) {
	evento := shared.NotificacaoEvent{
		ContaNum: contaNum,
		Mensagem: mensagem,
		Sucesso:  sucesso,
	}

	body, err := json.Marshal(evento)
	if err != nil {
		log.Printf("Erro ao serializar notificação: %v", err)
		return
	}

	err = canal.Publish(
		"",
		shared.FilaNotificacao,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Erro ao publicar notificação na fila: %v", err)
	}
}
