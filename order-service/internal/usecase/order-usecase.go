package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Nurda-zh/a1/order-service/internal/domain"
	"github.com/Nurda-zh/a1/order-service/internal/repository"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrStockInsufficient = errors.New("stock insufficient")
)

type OrderUsecase interface {
	CreateOrder(req *domain.CreateOrderRequest) (string, error)
	GetOrder(id string) (*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	ListOrdersByUser(userID string, page, pageSize int64) ([]*domain.Order, int64, error)
}

type orderUsecase struct {
	repo         repository.OrderRepo
	inventoryURL string
	httpClient   *http.Client
}

func NewOrderUsecase(r repository.OrderRepo, inventoryURL string) OrderUsecase {
	return &orderUsecase{
		repo:         r,
		inventoryURL: inventoryURL,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

type reserveItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type reserveReq struct {
	Items []reserveItem `json:"items"`
}

// CreateOrder: basic flow:
// 1. validate
// 2. attempt to reserve stock via Inventory API (POST /products/reserve or similar). For assignment we call /products/stock-decrement endpoint (mock).
// 3. if reserve ok, create order in DB and return id
func (u *orderUsecase) CreateOrder(req *domain.CreateOrderRequest) (string, error) {
	if req.UserID == "" {
		return "", errors.New("user_id required")
	}
	if len(req.Items) == 0 {
		return "", errors.New("items required")
	}

	// compute total (if price provided in items, use it; else 0)
	var total int64
	for _, it := range req.Items {
		if it.Quantity <= 0 {
			return "", fmt.Errorf("invalid quantity for product %s", it.ProductID)
		}
		total += int64(it.Quantity) * it.PriceCents
	}

	// reserve stock via inventory (synchronous). We expect inventory to expose POST /products/reserve
	reserve := reserveReq{}
	for _, it := range req.Items {
		reserve.Items = append(reserve.Items, reserveItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		})
	}
	body, _ := json.Marshal(reserve)
	reserveURL := fmt.Sprintf("%s/products/reserve", u.inventoryURL) // e.g. http://inventory:8001/api/products/reserve
	httpReq, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, reserveURL, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := u.httpClient.Do(httpReq)
	if err != nil {
		// inventory unreachable — for assignment we treat as failure
		return "", fmt.Errorf("inventory reserve failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// failed to reserve
		return "", ErrStockInsufficient
	}

	// build order entity
	o := &domain.Order{
		UserID:     req.UserID,
		Items:      req.Items,
		TotalCents: total,
		Status:     domain.StatusPending,
	}
	id, err := u.repo.Create(o)
	if err != nil {
		// optionally revert reserve — omitted for simplicity
		return "", err
	}
	return id, nil
}

func (u *orderUsecase) GetOrder(id string) (*domain.Order, error) {
	o, err := u.repo.GetByID(id)
	if err != nil {
		if err == repository.ErrOrderNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return o, nil
}

func (u *orderUsecase) UpdateStatus(id string, status domain.OrderStatus) error {
	// basic validation
	if status != domain.StatusCancelled && status != domain.StatusCompleted && status != domain.StatusPending {
		return errors.New("invalid status")
	}
	if err := u.repo.UpdateStatus(id, status); err != nil {
		if err == repository.ErrOrderNotFound {
			return ErrOrderNotFound
		}
		return err
	}
	return nil
}

func (u *orderUsecase) ListOrdersByUser(userID string, page, pageSize int64) ([]*domain.Order, int64, error) {
	return u.repo.ListByUser(userID, page, pageSize)
}
