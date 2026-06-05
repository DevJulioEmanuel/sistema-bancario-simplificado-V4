package publisher

import (
	"context"
	"encoding/json"
	"shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	canal *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	canal, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	filas := []string{
		shared.FilaDeposito,
		shared.FilaSaque,
		shared.FilaTransferencia,
		shared.FilaPagamento,
	}

	for _, fila := range filas {
		_, err = canal.QueueDeclare(fila, true, false, false, false, nil)
		if err != nil {
			return nil, err
		}
	}

	return &Publisher{canal: canal}, nil
}

func (p *Publisher) Publish(ctx context.Context, queueName string, transacao shared.TransacaoEvent) error {
	body, err := json.Marshal(transacao)

	if err != nil {
		return err
	}

	return p.canal.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
