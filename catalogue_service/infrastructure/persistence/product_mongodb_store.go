package persistence

import (
	"context"
	"github.com/tamararankovic/microservices_demo/catalogue_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "catalogue"
	COLLECTION = "product"
)

type ProductMongoDBStore struct {
	products *mongo.Collection
}

func NewProductMongoDBStore(client *mongo.Client) domain.ProductStore {
	products := client.Database(DATABASE).Collection(COLLECTION)
	return &ProductMongoDBStore{
		products: products,
	}
}

func (store *ProductMongoDBStore) Get(id primitive.ObjectID) (*domain.Product, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *ProductMongoDBStore) GetAll() ([]*domain.Product, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *ProductMongoDBStore) Insert(product *domain.Product) error {
	result, err := store.products.InsertOne(context.TODO(), product)
	if err != nil {
		return err
	}
	product.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *ProductMongoDBStore) DeleteAll() {
	store.products.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *ProductMongoDBStore) filter(filter interface{}) ([]*domain.Product, error) {
	cursor, err := store.products.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *ProductMongoDBStore) filterOne(filter interface{}) (product *domain.Product, err error) {
	result := store.products.FindOne(context.TODO(), filter)
	err = result.Decode(&product)
	return
}

func decode(cursor *mongo.Cursor) (products []*domain.Product, err error) {
	for cursor.Next(context.TODO()) {
		var product domain.Product
		err = cursor.Decode(&product)
		if err != nil {
			return
		}
		products = append(products, &product)
	}
	err = cursor.Err()
	return
}
