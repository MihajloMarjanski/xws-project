package handler_grpc

import (
	"requests-service/service"

	events "github.com/MihajloMarjanski/xws-project/common/saga/block_user"
	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
)

type BlockUserCommandHandler struct {
	requestService    *service.RequestsService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewBlockUserCommandHandler(requestService *service.RequestsService, publisher saga.Publisher, subscriber saga.Subscriber) (*BlockUserCommandHandler, error) {
	o := &BlockUserCommandHandler{
		requestService:    requestService,
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
	case events.UpdateConnectionRequest:
		success := handler.requestService.RemoveConnection(uint(command.Block.BlockId), uint(command.Block.UserId))
		if success {
			reply.Type = events.ConnectionRequestUpdated
		} else {
			reply.Type = events.ConnectionRequestNotUpdated
		}

	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
