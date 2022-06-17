package startup

import (
	cfg "api-gateway/startup/config"
	"context"
	"fmt"
	postGw "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
	requestGw "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	userGw "github.com/MihajloMarjanski/xws-project/common/proto/user_service"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"path/filepath"
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
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	creds, err6 := credentials.NewClientTLSFromFile("startup/certTLS/service.pem", "")
	if err6 != nil {
		log.Fatalf("could not process the credentials: %v", err6)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		panic(err)
	}
	requestEndpoint := fmt.Sprintf("%s:%s", server.config.RequestHost, server.config.RequestPort)
	err1 := requestGw.RegisterRequestsServiceHandlerFromEndpoint(context.TODO(), server.mux, requestEndpoint, opts)
	if err1 != nil {
		panic(err1)
	}
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err2 := postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err2 != nil {
		panic(err2)
	}
}

func (server *Server) Start() {
	crtPath, _ := filepath.Abs("./server.crt")
	keyPath, _ := filepath.Abs("./server.key")
	origins := handlers.AllowedOrigins([]string{"https://localhost:4300", "https://localhost:4300/**", "https://localhost:4300", "https://localhost:4300/**"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	headers := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Authorization", "Access-Control-Allow-Origin", "*"})

	log.Println("gateway started")
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), handlers.CORS(headers, methods, origins)(server.mux)))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", server.config.Port), crtPath, keyPath, handlers.CORS(headers, methods, origins)(server.mux)))
}
