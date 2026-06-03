package main

import (
	"banco-api/internal/handler"
	"banco-api/internal/publisher"
	"banco-api/internal/repository"
	"banco-api/internal/routes"
	"banco-api/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	pub, err := publisher.NewPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("erro ao conectar no RabbitMQ: %v", err)
	}
	bancoRepo := repository.NewBancoRepository()

	clienteRepo := repository.NewClienteRepository(
		bancoRepo,
	)

	contaRepo := repository.NewContaRepository(
		bancoRepo,
	)

	clienteService := service.NewClienteService(clienteRepo, contaRepo)
	contaService := service.NewContaService(contaRepo)

	clienteHandler := handler.NewClienteHandler(clienteService)
	contaHandler := handler.NewContaHandler(contaService, pub)

	r := gin.Default()

	routes.SetupRoutes(r, clienteHandler, contaHandler)

	r.Run(":8080")
}
