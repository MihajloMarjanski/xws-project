package startup

import (
	"github.com/tamararankovic/microservices_demo/catalogue_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var products = []*domain.Product{
	{
		Id:            getObjectId("623b0cc3a34d25d8567f9f82"),
		Name:          "name",
		ClothingBrand: "brand",
		Colors: []domain.Color{
			{
				Code: "R",
				Name: "Red",
			},
			{
				Code: "B",
				Name: "Blue",
			},
		},
	},
	{
		Id:            getObjectId("623b0cc3a34d25d8567f9f83"),
		Name:          "name2",
		ClothingBrand: "brand2",
		Colors: []domain.Color{
			{
				Code: "R",
				Name: "Red",
			},
			{
				Code: "G",
				Name: "Green",
			},
		},
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
