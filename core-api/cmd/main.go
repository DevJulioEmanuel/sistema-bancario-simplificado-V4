package main

import (
	"banco-api/internal/db"
	"banco-api/internal/handler"
	"banco-api/internal/publisher"
	"banco-api/internal/repository"
	"banco-api/internal/routes"
	"banco-api/internal/service"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	var pub *publisher.Publisher
	var err error
	for i := range 10 {
		pub, err = publisher.NewPublisher(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("aguardando RabbitMQ... tentativa %d/10", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("erro ao conectar no RabbitMQ: %v", err)
	}

	db, err := db.NewDB(dbPath)
	if err != nil {
		log.Fatalf("erro ao conectar no banco de dados: %v", err)
	}
	clienteRepo := repository.NewClienteRepository(
		db,
	)

	contaRepo := repository.NewContaRepository(
		db,
	)

	clienteService := service.NewClienteService(clienteRepo, contaRepo)
	contaService := service.NewContaService(contaRepo)

	clienteHandler := handler.NewClienteHandler(clienteService)
	contaHandler := handler.NewContaHandler(contaService, pub)

	r := gin.Default()

	routes.SetupRoutes(r, clienteHandler, contaHandler)

	r.Run(":8080")
}
