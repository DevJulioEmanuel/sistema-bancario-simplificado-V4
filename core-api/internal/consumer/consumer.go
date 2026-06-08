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
		log.Println("Consumidor de WebSocket escutando a fila_resposta_pagamentos...")

		for d := range msgs {
			log.Printf("[RabbitMQ] Mensagem recebida para WS: %s", string(d.Body))

			var resp RespostaFaturamento
			err := json.Unmarshal(d.Body, &resp)
			if err != nil {
				log.Printf("[WEB SOCKET] Erro ao decodificar JSON: %v", err)
				continue
			}

			log.Printf("[WEB SOCKET] Buscando WS para a conta [%d] com novo saldo [%.2f]", resp.NumeroConta, resp.NovoSaldo)

			handler.ClientesHub.Lock()
			conn, existe := handler.ClientesHub.Conexoes[resp.NumeroConta]

			if !existe {
				log.Printf("[WEB SOCKET] WS NÃO encontrado para a conta %d. Conexões ativas no Hub: %d", resp.NumeroConta, len(handler.ClientesHub.Conexoes))
			}
			handler.ClientesHub.Unlock()

			if existe {
				err = conn.WriteJSON(map[string]interface{}{
					"evento":     "saldo_atualizado",
					"novo_saldo": resp.NovoSaldo,
				})
				if err != nil {
					log.Printf("[WEB SOCKET] Erro ao enviar para o cliente via WS: %v", err)
				} else {
					log.Printf("[WEB SOCKET] Saldo enviado com sucesso para o cliente %d via WS!", resp.NumeroConta)
				}
			}
		}

	}()
}
