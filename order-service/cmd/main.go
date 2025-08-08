package main

import (
	"context"
	"log"
	"net/http"
	"time"

	config "github.com/Nurda-zh/a1/order-service/configs"
	"github.com/Nurda-zh/a1/order-service/internal/delivery/http/handler"
	infra "github.com/Nurda-zh/a1/order-service/internal/infra"
	"github.com/Nurda-zh/a1/order-service/internal/repository"
	"github.com/Nurda-zh/a1/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	client, err := infra.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	db := client.Database(cfg.MongoDBName)

	orderRepo := repository.NewMongoOrderRepo(db)
	orderUC := usecase.NewOrderUsecase(orderRepo, cfg.InventoryServiceURL)
	orderHandler := handler.NewOrderHandler(orderUC)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api")
	orderHandler.RegisterRoutes(api)

	addr := ":" + cfg.ListenPort
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("Order service running on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
