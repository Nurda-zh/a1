package repository

import (
	"context"
	"time"

	"github.com/Nurda-zh/a1/inventory-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	Update(ctx context.Context, id string, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]entity.Product, error)
}

type productRepository struct {
	col *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepository{
		col: db.Collection("products"),
	}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.col.InsertOne(ctx, product)
	return err
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var product entity.Product
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	return &product, err
}

func (r *productRepository) Update(ctx context.Context, id string, product *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": product})
	return err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *productRepository) List(ctx context.Context) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []entity.Product
	for cur.Next(ctx) {
		var p entity.Product
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
