package service

import (
	events "github.com/MihajloMarjanski/xws-project/common/saga/block_user"

	saga "github.com/MihajloMarjanski/xws-project/common/saga/messaging"
)

type BlockUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewBlockUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*BlockUserOrchestrator, error) {
	o := &BlockUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *BlockUserOrchestrator) Start(userId, blockId uint) error {
	event := &events.BlockUserCommand{
		Type: events.UpdateConnectionRequest,
		Block: events.BlockDetails{
			UserId:  uint64(userId),
			BlockId: uint64(blockId),
		},
	}

	return o.commandPublisher.Publish(event)
}

func (o *BlockUserOrchestrator) handle(reply *events.BlockUserReply) {
	command := events.BlockUserCommand{Block: reply.Block}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *BlockUserOrchestrator) nextCommandType(reply events.BlockUserReplyType) events.BlockUserCommandType {
	switch reply {
	case events.ConnectionRequestNotUpdated:
		return events.CancelBlock
	case events.ConnectionRequestUpdated:
		return events.BlockUser
	default:
		return events.UnknownCommand
	}
}
