package service

import (
	"fmt"
)

type ConnectionService struct {
}

func New() (*ConnectionService, error) {

	return &ConnectionService{}, nil
}

func (connectionService *ConnectionService) Connect(id1, id2 uint64) {
	t := fmt.Sprintf("CONNECTIG ID1: %[1]d and ID2: %[2]d\n", id1, id2)
	fmt.Println(t)
}

func (connectionService *ConnectionService) GetRecommendedConnections(id uint64) []uint64 {
	ids := []uint64{1, 2, 3}
	return ids
}
