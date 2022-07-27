package nats

import (
	//saga "common/saga/messaging"
	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn       *nats.EncodedConn
	subject    string
	queueGroup string
}

func NewNATSSubscriber(subject, queueGroup string) (saga.Subscriber, error) {
	conn, err := getConnection()
	if err != nil {
		return nil, err
	}
	encConn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		conn:       encConn,
		subject:    subject,
		queueGroup: queueGroup,
	}, nil
}

func (s *Subscriber) Subscribe(handler interface{}) error {
	_, err := s.conn.QueueSubscribe(s.subject, s.queueGroup, handler)
	if err != nil {
		return err
	}
	return nil
}
