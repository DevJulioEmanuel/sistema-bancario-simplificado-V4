package routes

import (
	"banco-api/internal/handler"
	"banco-api/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	clienteHandler *handler.ClienteHandler,
	contaHandler *handler.ContaHandler,
) {

	clientes := r.Group("/clientes")
	{
		clientes.POST("", clienteHandler.Cadastrar)
		clientes.POST("/login", clienteHandler.Login)
	}

	contas := r.Group("/contas/:num")

	contas.Use(middleware.JWTMiddleware())

	{
		contas.GET("/", contaHandler.ObterDados)
		contas.POST("/pagamento", contaHandler.Pagar)
		contas.POST("/depositar", contaHandler.Depositar)
		contas.POST("/sacar", contaHandler.Sacar)
		contas.POST("/transferir", contaHandler.Transferir)
		contas.GET("/extrato", contaHandler.Extrato)
		contas.GET("/rendimento/:meses", contaHandler.CalcularRendimento)
	}

	r.GET("/ws/:num", handler.HandleWebSocket)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
