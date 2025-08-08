package handler

import (
	"net/http"
	"strconv"

	"github.com/Nurda-zh/a1/order-service/internal/domain"
	"github.com/Nurda-zh/a1/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	uc usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) RegisterRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/orders")
	r.POST("", h.createOrder)
	r.GET("", h.listOrders)
	r.GET("/:id", h.getOrder)
	r.PATCH("/:id", h.patchOrder)
}

func (h *OrderHandler) createOrder(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.uc.CreateOrder(&req)
	if err != nil {
		if err == usecase.ErrStockInsufficient {
			c.JSON(http.StatusConflict, gin.H{"error": "stock insufficient"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", "/api/orders/"+id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *OrderHandler) getOrder(c *gin.Context) {
	id := c.Param("id")
	o, err := h.uc.GetOrder(id)
	if err != nil {
		if err == usecase.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

type patchStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *OrderHandler) patchOrder(c *gin.Context) {
	id := c.Param("id")
	var req patchStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := domain.OrderStatus(req.Status)
	if err := h.uc.UpdateStatus(id, status); err != nil {
		if err == usecase.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) listOrders(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 64)
	items, total, err := h.uc.ListOrdersByUser(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
