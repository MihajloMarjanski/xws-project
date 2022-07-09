package service

type ConnectionService struct {
}

func New() (*ConnectionService, error) {

	return &ConnectionService{}, nil
}

func (connectionService *ConnectionService) Connect(id1, id2 int64) {

}
