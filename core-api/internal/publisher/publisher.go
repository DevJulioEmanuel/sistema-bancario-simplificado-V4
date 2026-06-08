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

	_, err = canal.QueueDeclare(shared.FilaTransacoes, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Publisher{canal: canal}, nil
}

func (p *Publisher) Publish(ctx context.Context, transacao shared.TransacaoEvent) error {
	body, err := json.Marshal(transacao)
	if err != nil {
		return err
	}

	return p.canal.PublishWithContext(
		ctx,
		"",
		shared.FilaTransacoes,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Channel() *amqp.Channel {
	return p.canal
}
