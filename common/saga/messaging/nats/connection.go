package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func getConnection() (*nats.Conn, error) {
	url := fmt.Sprintf("nats://ruser:T0pS3cr3t@nats:4222")
	connection, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
