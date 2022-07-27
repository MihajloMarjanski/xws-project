package handler_grpc

import (
	"user-service/service"

	events "github.com/MihajloMarjanski/xws-project/common/saga/block_user"
	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
)

type BlockUserCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewBlockUserCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*BlockUserCommandHandler, error) {
	o := &BlockUserCommandHandler{
		userService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *BlockUserCommandHandler) handle(command *events.BlockUserCommand) {

	reply := events.BlockUserReply{Block: command.Block}

	switch command.Type {
	case events.BlockUser:
		reply.Type = events.UserBlocked
	case events.CancelBlock:
		reply.Type = events.UserBlockCancelled
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
