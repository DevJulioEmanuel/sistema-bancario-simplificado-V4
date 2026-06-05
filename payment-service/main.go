package main

import (
	"encoding/json"
	"log"
	"os"
	"payment-service/internal/processor"
	"shared"
	"time"

	"github.com/glebarez/sqlite"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./banco.db"
	}

	db, errdb := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if errdb != nil {
		log.Fatalf("erro ao conectar no banco: %v", errdb)
	}
	proc := processor.NewProcessor(db)

	var conn *amqp.Connection
	var err error
	for i := range 10 {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("aguardando RabbitMQ... tentativa %d/10", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("erro ao conectar no RabbitMQ: %v", err)
	}

	canal, err := conn.Channel()
	if err != nil {
		log.Fatalf("erro ao abrir canal: %v", err)
	}
	defer canal.Close()

	filas := []string{"deposito", "saque", "transferencia", "pagamento"}
	for _, fila := range filas {
		_, err = canal.QueueDeclare(fila, true, false, false, false, nil)
		if err != nil {
			log.Fatalf("erro ao declarar fila %s: %v", fila, err)
		}
	}

	deposito, err := canal.Consume("deposito", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("erro: %v", err)
	}

	saque, err := canal.Consume("saque", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("erro: %v", err)
	}

	transferencia, err := canal.Consume("transferencia", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("erro: %v", err)
	}

	pagamento, err := canal.Consume("pagamento", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("erro: %v", err)
	}

	for {
		select {
		case msg := <-deposito:
			var transacao shared.TransacaoEvent
			err := json.Unmarshal(msg.Body, &transacao)
			if err != nil {
				log.Printf("erro ao deserializar: %v", err)
				msg.Nack(false, true) // recoloca na fila
				continue
			}
			err = proc.ProcessarDeposito(transacao)
			if err != nil {
				log.Printf("erro ao processar depósito: %v", err)
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		case msg := <-saque:
			var transacao shared.TransacaoEvent
			err := json.Unmarshal(msg.Body, &transacao)
			if err != nil {
				log.Printf("erro ao deserializar: %v", err)
				msg.Nack(false, true)
				continue
			}
			err = proc.ProcessarSaque(transacao)
			if err != nil {
				log.Printf("erro ao processar saque: %v", err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		case msg := <-transferencia:
			var transacao shared.TransacaoEvent
			err := json.Unmarshal(msg.Body, &transacao)
			if err != nil {
				log.Printf("erro ao deserializar: %v", err)
				msg.Nack(false, true)
				continue
			}
			err = proc.ProcessarTransferencia(transacao)
			if err != nil {
				log.Printf("erro ao processar transferência: %v", err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		case msg := <-pagamento:
			var transacao shared.TransacaoEvent
			err := json.Unmarshal(msg.Body, &transacao)
			if err != nil {
				log.Printf("erro ao deserializar: %v", err)
				msg.Nack(false, true)
				continue
			}
			err = proc.ProcessarPagamento(transacao)
			if err != nil {
				log.Printf("erro ao processar pagamento: %v", err)
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)
		}
	}
}
