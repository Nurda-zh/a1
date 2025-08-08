package handler

import (
	"github.com/Nurda-zh/a1/inventory-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, ph *handler.ProductHandler) {
	r.POST("/products", ph.CreateProduct)
	r.GET("/products/:id", ph.GetProduct)
	r.PATCH("/products/:id", ph.UpdateProduct)
	r.DELETE("/products/:id", ph.DeleteProduct)
	r.GET("/products", ph.ListProducts)
}
