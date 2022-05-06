package api

import (
	"context"
	"github.com/tamararankovic/microservices_demo/catalogue_service/application"
	pb "github.com/tamararankovic/microservices_demo/common/proto/catalogue_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	pb.UnimplementedCatalogueServiceServer
	service *application.ProductService
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (handler *ProductHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	product, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	productPb := mapProduct(product)
	response := &pb.GetResponse{
		Product: productPb,
	}
	return response, nil
}

func (handler *ProductHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	products, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Products: []*pb.Product{},
	}
	for _, product := range products {
		current := mapProduct(product)
		response.Products = append(response.Products, current)
	}
	return response, nil
}
