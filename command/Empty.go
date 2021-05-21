package command

import (
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type EmptyCommandHandler struct{}

func (EmptyCommandHandler) Handle(frame *CommandFrame) error {
	typ, err := types.Deserialize(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetById(frame.ClientId)
	empty, err := cl.Empty(typ)
	if err != nil {
		return err
	}
	binEmpty := []byte{0}
	if empty {
		binEmpty[0] = 1
	}
	return cl.SendToClient(binEmpty)
}
