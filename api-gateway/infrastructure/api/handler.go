package api

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

type Handler interface {
	Init(mux *runtime.ServeMux)
}
