package command

import (
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type EmptyCommandHandler struct{}

func (EmptyCommandHandler) Handle(frame CommandFrame) error {
	typ, err := types.Deserialize(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetById(frame.ClientId)
	_, err = cl.Empty(typ)
	return err
}

type emptyCommandContent struct {
	typ types.Type
}
