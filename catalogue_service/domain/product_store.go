package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductStore interface {
	Get(id primitive.ObjectID) (*Product, error)
	GetAll() ([]*Product, error)
	Insert(product *Product) error
	DeleteAll()
}
