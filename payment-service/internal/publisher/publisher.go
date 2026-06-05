package publisher

import (
	"encoding/json"
    "log"
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