package command

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/google/uuid"
	"log"
)

type Handler interface {
	Handle(frame *CommandFrame) error
}

type DefaultHandler struct{}

func (DefaultHandler) Handle(*CommandFrame) error {
	return errors.New("unsupported command")
}

type CommandFrame struct {
	ClientId  uuid.UUID
	CommandId uint8
	Size      uint64
	Data      []byte
}

func getClientId(rawFrame []byte) (uuid.UUID, error) {
	if len(rawFrame) < 25 {
		return uuid.Nil, errors.New("insufficient input length")
	}
	rawClientId := rawFrame[0:16]
	return uuid.FromBytes(rawClientId)
}

func parseCommandFrame(rawFrame []byte) (CommandFrame, error) {
	if len(rawFrame) < 25 {
		return CommandFrame{}, errors.New("insufficient input length")
	}
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

	clientId, err := getClientId(rawFrame)
	if err != nil {
		return CommandFrame{}, err
	}

	return CommandFrame{
		ClientId:  clientId,
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

func Submit(rawFrame []byte) error {
	clientId, err := getClientId(rawFrame) // For faster locking
	if err != nil {
		return err
	}
	cl := client.Clients.GetById(clientId)
	if cl != nil {
		mutex := cl.GetCommandMutex()
		mutex.Lock()
		defer mutex.Unlock()
	}

	frame, err := parseCommandFrame(rawFrame)
	if err != nil {
		return err
	}

	return frame.Handle()
}

type handlerError struct {
	frame *CommandFrame
	cause error
}

func (e *handlerError) Error() string {
	return fmt.Sprintf("Command %d from client %s failed: %s", e.frame.CommandId, e.frame.ClientId, e.cause.Error())
}

func (frame *CommandFrame) Handle() error {
	log.Printf("Client %s issued command %d", frame.ClientId.String(), frame.CommandId)
	var handler Handler
	switch frame.CommandId {
	case RegisterCommandId:
		handler = RegisterCommandHandler{}
	case AcceptTypeCommandId:
		handler = AcceptTypeCommandHandler{}
	case SendCommandId:
		handler = SendCommandHandler{}
	case GetCommandId:
		handler = GetCommandHandler{}
	case EmptyCommandId:
		handler = EmptyCommandHandler{}
	default:
		handler = DefaultHandler{}
	}
	err := handler.Handle(frame)
	if err == nil {
		return nil
	}
	return &handlerError{
		frame: frame,
		cause: err,
	}
}
