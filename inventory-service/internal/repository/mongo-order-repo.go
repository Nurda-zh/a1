package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nurda-zh/a1/order-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrOrderNotFound = errors.New("order not found")

type MongoOrderRepo struct {
	coll *mongo.Collection
}

func NewMongoOrderRepo(db *mongo.Database) *MongoOrderRepo {
	return &MongoOrderRepo{
		coll: db.Collection("orders"),
	}
}

func (r *MongoOrderRepo) Create(order *domain.Order) (string, error) {
	now := time.Now().UTC()
	order.CreatedAt = now
	order.UpdatedAt = now
	if order.Status == "" {
		order.Status = domain.StatusPending
	}
	oid := primitive.NewObjectID()
	doc := bson.M{
		"_id":         oid,
		"user_id":     order.UserID,
		"items":       order.Items,
		"total_cents": order.TotalCents,
		"status":      order.Status,
		"created_at":  order.CreatedAt,
		"updated_at":  order.UpdatedAt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	return oid.Hex(), nil
}

func (r *MongoOrderRepo) GetByID(id string) (*domain.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res bson.M
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	o := &domain.Order{}
	if idv, ok := res["_id"].(primitive.ObjectID); ok {
		o.ID = idv.Hex()
	} else {
		o.ID = fmt.Sprintf("%v", res["_id"])
	}
	if v, ok := res["user_id"].(string); ok {
		o.UserID = v
	}
	if items, ok := res["items"].(primitive.A); ok {
		// decode using bson -> marshal/unmarshal is simpler but to avoid import, iterate
		var parsed []domain.OrderItem
		for _, it := range items {
			if m, ok := it.(bson.M); ok {
				parsed = append(parsed, domain.OrderItem{
					ProductID:  getString(m["product_id"]),
					Quantity:   getInt(m["quantity"]),
					PriceCents: getInt64(m["price_cents"]),
				})
			}
		}
		o.Items = parsed
	} else if items, ok := res["items"].([]interface{}); ok {
		var parsed []domain.OrderItem
		for _, it := range items {
			if m, ok := it.(bson.M); ok {
				parsed = append(parsed, domain.OrderItem{
					ProductID:  getString(m["product_id"]),
					Quantity:   getInt(m["quantity"]),
					PriceCents: getInt64(m["price_cents"]),
				})
			}
		}
		o.Items = parsed
	}
	o.TotalCents = getInt64(res["total_cents"])
	if s, ok := res["status"].(string); ok {
		o.Status = domain.OrderStatus(s)
	}
	if t, ok := res["created_at"].(primitive.DateTime); ok {
		o.CreatedAt = t.Time()
	}
	if t, ok := res["updated_at"].(primitive.DateTime); ok {
		o.UpdatedAt = t.Time()
	}
	return o, nil
}

func (r *MongoOrderRepo) UpdateStatus(id string, status domain.OrderStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrOrderNotFound
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"status": status, "updated_at": time.Now().UTC()},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (r *MongoOrderRepo) ListByUser(userID string, page, pageSize int64) ([]*domain.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	f := bson.M{"user_id": userID}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	total, err := r.coll.CountDocuments(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip((page - 1) * pageSize).SetLimit(pageSize).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cur, err := r.coll.Find(ctx, f, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	var out []*domain.Order
	for cur.Next(ctx) {
		var res bson.M
		if err := cur.Decode(&res); err != nil {
			return nil, 0, err
		}
		o := &domain.Order{}
		if idv, ok := res["_id"].(primitive.ObjectID); ok {
			o.ID = idv.Hex()
		} else {
			o.ID = fmt.Sprintf("%v", res["_id"])
		}
		if v, ok := res["user_id"].(string); ok {
			o.UserID = v
		}
		o.TotalCents = getInt64(res["total_cents"])
		if s, ok := res["status"].(string); ok {
			o.Status = domain.OrderStatus(s)
		}
		if t, ok := res["created_at"].(primitive.DateTime); ok {
			o.CreatedAt = t.Time()
		}
		if t, ok := res["updated_at"].(primitive.DateTime); ok {
			o.UpdatedAt = t.Time()
		}
		// items parse (similar to GetByID)
		if items, ok := res["items"].(primitive.A); ok {
			var parsed []domain.OrderItem
			for _, it := range items {
				if m, ok := it.(bson.M); ok {
					parsed = append(parsed, domain.OrderItem{
						ProductID:  getString(m["product_id"]),
						Quantity:   getInt(m["quantity"]),
						PriceCents: getInt64(m["price_cents"]),
					})
				}
			}
			o.Items = parsed
		}
		out = append(out, o)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// helpers (same as inventory repo utils)
func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float64:
		return int(t)
	default:
		return 0
	}
}

func getInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case float64:
		return int64(t)
	default:
		return 0
	}
}
