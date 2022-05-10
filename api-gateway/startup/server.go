package startup

import (
	cfg "api-gateway/startup/config"
	"context"
	"fmt"
	requestGw "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	userGw "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	return server
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		panic(err)
	}
	requestEndpoint := fmt.Sprintf("%s:%s", server.config.RequestHost, server.config.RequestPort)
	err1 := requestGw.RegisterRequestsServiceHandlerFromEndpoint(context.TODO(), server.mux, requestEndpoint, opts)
	if err1 != nil {
		panic(err)
	}

}

func (server *Server) Start() {
	log.Println("gateway started")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))

}
