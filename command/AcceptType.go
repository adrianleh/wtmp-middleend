package command

import (
	"errors"
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type AcceptTypeCommandHandler struct{}

func (AcceptTypeCommandHandler) Handle(frame *CommandFrame) error {
	typ, err := types.Deserialize(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetById(frame.ClientId)
	if cl == nil {
		return errors.New("client not found")
	}
	err = cl.RegisterType(typ)
	ret := []byte{0}
	if err != nil {
		ret[0] = 1
	}
	sendErr := cl.SendToClient(ret)
	if err == nil && sendErr != nil {
		return sendErr
	}
	return err
}
