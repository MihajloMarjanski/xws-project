package api

import (
	"github.com/tamararankovic/microservices_demo/catalogue_service/domain"
	pb "github.com/tamararankovic/microservices_demo/common/proto/catalogue_service"
)

func mapProduct(product *domain.Product) *pb.Product {
	productPb := &pb.Product{
		Id:            product.Id.Hex(),
		Name:          product.Name,
		ClothingBrand: product.ClothingBrand,
	}
	for _, color := range product.Colors {
		productPb.Colors = append(productPb.Colors, &pb.Color{
			Code: color.Code,
			Name: color.Name,
		})
	}
	return productPb
}
