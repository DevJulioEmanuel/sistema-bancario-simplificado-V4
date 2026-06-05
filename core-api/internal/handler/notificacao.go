package handler

import (
	"banco-api/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificacaoHandler struct {
	db *gorm.DB
}

func NewNotificacaoHandler(db *gorm.DB) *NotificacaoHandler {
	return &NotificacaoHandler{db: db}
}

func (h *NotificacaoHandler) GetNotificacoes(c *gin.Context) {
	contaNumStr := c.Param("num")
	contaNum, err := strconv.Atoi(contaNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Número de conta inválido"})
		return
	}

	var notificacoes []model.Notificacao

	result := h.db.Where("conta_num = ?", contaNum).
		Order("created_at desc").
		Limit(20).
		Find(&notificacoes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar notificações"})
		return
	}

	c.JSON(http.StatusOK, notificacoes)
}
