package command

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
)

type Handler interface {
	Handle(frame CommandFrame) error
}

type DefaultHandler struct{}

func (DefaultHandler) Handle(frame CommandFrame) error {
	return errors.New("unsupported command")
}

type CommandFrame struct {
	ClientId  uuid.UUID
	CommandId uint8
	Size      uint64
	Data      []byte
}

func ParseCommandFrame(rawFrame []byte) (CommandFrame, error) {
	if len(rawFrame) < 25 {
		return CommandFrame{}, errors.New("insufficient input length")
	}
	uuidRaw := rawFrame[0:16]
	commandIdRaw := rawFrame[16]
	sizeRaw := rawFrame[17:25]

	size := binary.BigEndian.Uint64(sizeRaw)

	if uint64(len(rawFrame)-25) != size {
		return CommandFrame{}, errors.New("data size mismatch")
	}

	var data []byte
	if size == 0 {
		data = make([]byte, 0)
	} else {
		data = rawFrame[25:]
	}

	commandId, err := uuid.FromBytes(uuidRaw)
	if err != nil {
		return CommandFrame{}, err
	}

	return CommandFrame{
		ClientId:  commandId,
		CommandId: commandIdRaw,
		Size:      size,
		Data:      data,
	}, nil
}

const (
	RegisterCommandId        = uint8(0)
	RegisterSubTypeCommandId = uint8(1)
	AcceptTypeCommandId      = uint8(2)
	SendCommandId            = uint8(3)
	GetCommandId             = uint8(4)
	EmptyCommandId           = uint8(5)
)

func Handle(frame CommandFrame) error {
	var handler Handler
	switch frame.CommandId {
	case RegisterCommandId:
		handler = RegisterCommandHandler{}
	default:
		handler = DefaultHandler{}
	}
	return handler.Handle(frame)
}
