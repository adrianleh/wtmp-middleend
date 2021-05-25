package command

import (
	"errors"
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type GetCommandHandler struct{}

func (GetCommandHandler) Handle(frame *CommandFrame) error {
	typ, err := types.Deserialize(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetById(frame.ClientId)
	if cl == nil {
		return errors.New("client not found")
	}
	data, err := cl.Pop(typ)
	if err != nil {
		return err
	}
	return cl.SendToClient(data)
}
