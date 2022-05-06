package application

import (
	"github.com/tamararankovic/microservices_demo/catalogue_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	store domain.ProductStore
}

func NewProductService(store domain.ProductStore) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (service *ProductService) Get(id primitive.ObjectID) (*domain.Product, error) {
	return service.store.Get(id)
}

func (service *ProductService) GetAll() ([]*domain.Product, error) {
	return service.store.GetAll()
}
