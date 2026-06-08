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
	defer conn.Close()

	canal, err := conn.Channel()
	if err != nil {
		log.Fatalf("erro ao abrir canal: %v", err)
	}
	defer canal.Close()

	err = canal.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("erro ao configurar QoS: %v", err)
	}

	_, err = canal.QueueDeclare(shared.FilaTransacoes, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("erro ao declarar fila: %v", err)
	}

	msgs, err := canal.Consume(shared.FilaTransacoes, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("erro ao consumir mensagens: %v", err)
	}

	log.Println("Serviço iniciado. Aguardando transações em ordem restrita...")

	for msg := range msgs {
		var transacao shared.TransacaoEvent
		err := json.Unmarshal(msg.Body, &transacao)

		if err != nil {
			log.Printf("erro ao deserializar payload: %v", err)
			msg.Nack(false, false)
			continue
		}

		var processErr error
		switch transacao.Tipo {
		case "deposito":
			processErr = proc.ProcessarDeposito(canal, transacao)
		case "saque":
			processErr = proc.ProcessarSaque(canal, transacao)
		case "transferencia":
			processErr = proc.ProcessarTransferencia(canal, transacao)
		case "pagamento":
			processErr = proc.ProcessarPagamento(canal, transacao)
		default:
			log.Printf("tipo de operação desconhecida: %s", transacao.Tipo)
			msg.Nack(false, false)
			continue
		}

		if processErr != nil {
			log.Printf("erro ao processar %s: %v", transacao.Tipo, processErr)
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}
}
