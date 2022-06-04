package block_user

type BlockDetails struct {
	UserId  uint64
	BlockId uint64
}

type BlockUserCommandType int8

const (
	UpdateConnectionRequest BlockUserCommandType = iota
	BlockUser
	CancelBlock
	UnknownCommand
)

type BlockUserCommand struct {
	Block BlockDetails
	Type  BlockUserCommandType
}

type BlockUserReplyType int8

const (
	ConnectionRequestUpdated BlockUserReplyType = iota
	ConnectionRequestNotUpdated
	UserBlocked
	UserBlockCancelled
	UnknownReply
)

type BlockUserReply struct {
	Block BlockDetails
	Type  BlockUserReplyType
}
