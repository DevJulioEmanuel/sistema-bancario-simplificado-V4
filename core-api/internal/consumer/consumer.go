package consumer

import (
	"banco-api/internal/handler"
	"encoding/json"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RespostaFaturamento struct {
	NumeroConta int     `json:"numero_conta"`
	NovoSaldo   float64 `json:"novo_saldo"`
}

func IniciarConsumidorResposta(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"fila_resposta_pagamentos",
		true,                       
		false,                      
		false,                      
		false,                      
		nil,             
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila de resposta: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao registrar consumidor de resposta: %v", err)
	}

	go func() {
		for d := range msgs {
			var resp RespostaFaturamento
			err := json.Unmarshal(d.Body, &resp)
			if err != nil {
				log.Printf("Erro ao decodificar JSON de resposta: %v", err)
				continue
			}

			handler.ClientesHub.Lock()
			conn, existe := handler.ClientesHub.Conexoes[resp.NumeroConta]
			handler.ClientesHub.Unlock()
			if existe {
				err = conn.WriteJSON(map[string]interface{}{
					"evento":     "saldo_atualizado",
					"novo_saldo": resp.NovoSaldo,
				})
				if err != nil {
					log.Printf("Erro ao enviar dados via WebSocket: %v", err)
				}
			}
		}
	}()
}