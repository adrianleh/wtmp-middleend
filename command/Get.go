package command

import (
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type GetCommandHandler struct{}

func (GetCommandHandler) Handle(frame CommandFrame) error {
	typ, err := types.Deserialize(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetById(frame.ClientId)
	_, err = cl.Pop(typ)
	return err
}

type getCommandContent struct {
	typ types.Type
}
