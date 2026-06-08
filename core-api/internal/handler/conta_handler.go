package handler

import (
	dtorequest "banco-api/internal/dto/request"
	dtoresponse "banco-api/internal/dto/response"
	"banco-api/internal/publisher"
	"banco-api/internal/service"
	"net/http"
	"shared"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContaHandler struct {
	service   *service.ContaService
	publisher *publisher.Publisher
}

func NewContaHandler(service *service.ContaService, publisher *publisher.Publisher) *ContaHandler {
	return &ContaHandler{
		service:   service,
		publisher: publisher,
	}
}

func (h *ContaHandler) ObterDados(c *gin.Context) {
	numero, _ := strconv.Atoi(c.Param("num"))

	conta, err := h.service.BuscarConta(numero)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	tipo := 1

	if conta.Tipo == 2 {
		tipo = 2
	}

	response := dtoresponse.ContaResponse{
		Numero:  conta.Numero,
		Titular: conta.Titular.Nome,
		Saldo:   conta.Saldo,
		Tipo:    tipo,
		Limite:  1200,
		Taxa:    0.005,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ContaHandler) Pagar(c *gin.Context) {

	numero, err := strconv.Atoi(c.Param("num"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "número inválido",
		})
		return
	}

	var req dtorequest.PagamentoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": err.Error(),
		})
		return
	}

	err = h.publisher.Publish(c, shared.TransacaoEvent{
		Tipo:            "pagamento",
		NumeroConta:     numero,
		Valor:           req.Valor,
		NumContaDestino: 0,
		Descricao:       req.Descricao,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"mensagem": "pagamento enviado para processamento"})
}

func (h *ContaHandler) Depositar(c *gin.Context) {
	numero, _ := strconv.Atoi(c.Param("num"))

	var req dtorequest.ValorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	err := h.publisher.Publish(c, shared.TransacaoEvent{
		Tipo:            "deposito",
		NumeroConta:     numero,
		Valor:           req.Valor,
		NumContaDestino: 0,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"mensagem": "depósito enviado para processamento"})
}

func (h *ContaHandler) Sacar(c *gin.Context) {
	numero, _ := strconv.Atoi(c.Param("num"))

	var req dtorequest.ValorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	err := h.publisher.Publish(c, shared.TransacaoEvent{
		Tipo:            "saque",
		NumeroConta:     numero,
		Valor:           req.Valor,
		NumContaDestino: 0,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"mensagem": "saque enviado para processamento"})
}

func (h *ContaHandler) Transferir(c *gin.Context) {
	numero, _ := strconv.Atoi(c.Param("num"))

	var req dtorequest.TransferenciaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	err := h.publisher.Publish(c, shared.TransacaoEvent{
		Tipo:            "transferencia",
		NumeroConta:     numero,
		Valor:           req.Valor,
		NumContaDestino: req.NumDestino,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"mensagem": "transferência enviada para processamento"})
}

func (h *ContaHandler) Extrato(c *gin.Context) {
	numero, err := strconv.Atoi(c.Param("num"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "número inválido",
		})
		return
	}

	conta, err := h.service.ObterExtrato(numero)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": err.Error(),
		})
		return
	}

	var transacoes []dtoresponse.TransacaoResponse

	for _, t := range conta.Historico {
		transacoes = append(transacoes, dtoresponse.TransacaoResponse{
			Tipo:      t.Tipo,
			Descricao: t.Descricao,
			Valor:     t.Valor,
			Data:      t.Data.Format("02/01/2006 15:04:05"),
		})
	}

	response := dtoresponse.ExtratoResponse{
		NumeroConta: conta.Numero,
		Titular:     conta.Titular.Nome,
		SaldoAtual:  conta.Saldo,
		Transacoes:  transacoes,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ContaHandler) CalcularRendimento(c *gin.Context) {

	numero, err := strconv.Atoi(c.Param("num"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "número inválido",
		})
		return
	}

	meses, err := strconv.Atoi(c.Param("meses"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "quantidade de meses inválida",
		})
		return
	}

	saldoProjetado, rendimento, conta, err :=
		h.service.CalcularRendimento(numero, meses)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": err.Error(),
		})
		return
	}

	response := dtoresponse.RendimentoResponse{
		NumeroConta:    conta.Numero,
		SaldoAtual:     conta.Saldo,
		Meses:          meses,
		TaxaMensal:     0.005,
		SaldoProjetado: saldoProjetado,
		Rendimento:     rendimento,
	}

	c.JSON(http.StatusOK, response)
}
