package main

import (
	"context"
	"log"
	"time"

	config "github.com/Nurda-zh/a1/order-service/configs"
	"github.com/Nurda-zh/a1/order-service/internal/delivery/http/handler"
	infra "github.com/Nurda-zh/a1/order-service/internal/infra"
	"github.com/Nurda-zh/a1/order-service/internal/repository"
	"github.com/Nurda-zh/a1/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	client, err := infra.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	db := client.Database(cfg.Database)

	orderRepo := repository.NewMongoOrderRepo(db)
	orderUC := usecase.NewOrderUsecase(orderRepo, cfg.InventoryServiceURL)
	orderHandler := handler.NewOrderHandler(orderUC)

	r := gin.Default()
	api := r.Group("/api")
	orderHandler.RegisterRoutes(api)

	log.Printf("Order service running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
