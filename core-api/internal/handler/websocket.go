package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Hub struct {
	sync.RWMutex
	Conexoes map[int]*websocket.Conn
}

var ClientesHub = &Hub{
	Conexoes: make(map[int]*websocket.Conn),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	numContaStr := c.Param("num")
	
	numeroConta, err := strconv.Atoi(numContaStr)
	if err != nil {
		log.Printf("Número de conta inválido no WebSocket: %s", numContaStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "número de conta inválido"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Erro ao fazer upgrade para WebSocket: %v", err)
		return
	}

	ClientesHub.Lock()
	ClientesHub.Conexoes[numeroConta] = conn
	ClientesHub.Unlock()

	log.Printf("[WebSocket] Cliente da conta %d conectado com sucesso!", numeroConta)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WebSocket] Cliente da conta %d desconectou.", numeroConta)

			ClientesHub.Lock()
			delete(ClientesHub.Conexoes, numeroConta)
			ClientesHub.Unlock()
			
			conn.Close()
			break
		}
	}
}