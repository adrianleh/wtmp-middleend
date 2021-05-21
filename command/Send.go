package command

import (
	"encoding/binary"
	"errors"
	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type SendCommandHandler struct{}

func (SendCommandHandler) Handle(frame *CommandFrame) error {
	content, err := sendData(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetByName(content.target)
	if cl == nil {
		return errors.New("client not found")
	}
	return cl.Push(content.typ, content.msg)
}

type sendCommandContent struct {
	typ    types.Type
	target string
	msg    []byte
}

func sendData(data []byte) (sendCommandContent, error) {
	if len(data) < 8 {
		return sendCommandContent{}, errors.New("data must at least have delimiters")
	}
	nameLen := binary.BigEndian.Uint32(data[0:4])
	typeLen := binary.BigEndian.Uint32(data[4:8])

	nameStartIdx := uint32(8)
	nameEndIdx := nameStartIdx + nameLen
	typeStartIdx := nameEndIdx
	typeEndIdx := typeStartIdx + typeLen

	if uint32(len(data)) < typeEndIdx {
		return sendCommandContent{}, errors.New("data too short")
	}

	nameRaw := data[nameStartIdx:nameEndIdx]
	typeRaw := data[typeStartIdx:typeEndIdx]

	name := string(nameRaw)

	typ, err := types.Deserialize(typeRaw)
	if err != nil {
		return sendCommandContent{}, err
	}

	msgStartIdx := typeEndIdx

	var msg []byte
	if uint32(len(data)) == msgStartIdx {
		msg = make([]byte, 0)
	} else {
		msg = data[msgStartIdx:]
	}

	return sendCommandContent{
		typ:    typ,
		target: name,
		msg:    msg,
	}, nil
}
