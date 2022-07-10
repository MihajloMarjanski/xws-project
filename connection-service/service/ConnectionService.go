package service

import (
	"connection-service/model"
	"connection-service/repo"
	"fmt"
	"strconv"
)

type ConnectionService struct {
	repository *repo.ConnectionRepository
}

func New() (*ConnectionService, error) {
	repository := repo.NewConnectionRepository()
	return &ConnectionService{repository: repository}, nil
}

func (connectionService *ConnectionService) Connect(id1, id2 uint64) {
	t := fmt.Sprintf("CONNECTIG ID1: %[1]d and ID2: %[2]d\n", id1, id2)
	fmt.Println(t)
	connection := model.Connection{UserOne: uint(id1), UserTwo: uint(id2)}
	connectionService.repository.CreateUserConnection(connection)

}

func (connectionService *ConnectionService) GetRecommendedConnections(id uint64) []uint64 {
	ids := []uint64{1, 2, 3}
	user := model.User{UserId: strconv.FormatUint(uint64(id), 10)}
	ids, _ = connectionService.repository.FindRecommendationsForUser(user)
	return ids
}
