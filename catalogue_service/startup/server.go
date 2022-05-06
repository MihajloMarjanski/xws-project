package startup

import (
	"fmt"
	"github.com/tamararankovic/microservices_demo/catalogue_service/application"
	"github.com/tamararankovic/microservices_demo/catalogue_service/domain"
	"github.com/tamararankovic/microservices_demo/catalogue_service/infrastructure/api"
	"github.com/tamararankovic/microservices_demo/catalogue_service/infrastructure/persistence"
	"github.com/tamararankovic/microservices_demo/catalogue_service/startup/config"
	catalogue "github.com/tamararankovic/microservices_demo/common/proto/catalogue_service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	productStore := server.initProductStore(mongoClient)

	productService := server.initProductService(productStore)

	productHandler := server.initProductHandler(productService)

	server.startGrpcServer(productHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.CatalogueDBHost, server.config.CatalogueDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initProductStore(client *mongo.Client) domain.ProductStore {
	store := persistence.NewProductMongoDBStore(client)
	store.DeleteAll()
	for _, product := range products {
		err := store.Insert(product)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initProductService(store domain.ProductStore) *application.ProductService {
	return application.NewProductService(store)
}

func (server *Server) initProductHandler(service *application.ProductService) *api.ProductHandler {
	return api.NewProductHandler(service)
}

func (server *Server) startGrpcServer(productHandler *api.ProductHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	catalogue.RegisterCatalogueServiceServer(grpcServer, productHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
