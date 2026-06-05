package consumer

import (
	"banco-api/internal/model"
	"encoding/json"
	"log"
	"shared"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func IniciarConsumidorNotificacoes(canal *amqp091.Channel, db *gorm.DB) {
	_, err := canal.QueueDeclare(
		shared.FilaNotificacao,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila de notificações: %v", err)
	}

	mensagens, err := canal.Consume(
		shared.FilaNotificacao,
		"core-api-notificacoes",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao registrar consumidor de notificações: %v", err)
	}

	go func() {
		for d := range mensagens {
			var evento shared.NotificacaoEvent
			if err := json.Unmarshal(d.Body, &evento); err != nil {
				log.Printf("Erro ao decodificar notificação: %v", err)
				continue
			}

			notificacao := model.Notificacao{
				ContaNum: evento.ContaNum,
				Mensagem: evento.Mensagem,
				Sucesso:  evento.Sucesso,
			}

			if err := db.Create(&notificacao).Error; err != nil {
				log.Printf("Erro ao salvar notificação no banco: %v", err)
			} else {
				log.Printf("Notificação salva para a conta %d", evento.ContaNum)
			}
		}
	}()
}
