module connection-service

go 1.18

replace github.com/MihajloMarjanski/xws-project/common => ../common

require (
	github.com/MihajloMarjanski/xws-project/common v0.0.0-00010101000000-000000000000
	github.com/neo4j/neo4j-go-driver/v4 v4.4.3
	google.golang.org/grpc v1.47.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.9.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
