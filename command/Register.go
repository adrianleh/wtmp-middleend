package command

import (
	"encoding/binary"
	"errors"
	"github.com/adrianleh/WTMP-middleend/client"
)

type RegisterCommandHandler struct{}

func (RegisterCommandHandler) Handle(frame CommandFrame) error {
	content, err := parseData(frame.Data)
	if err != nil {
		return err
	}

	cl, err := client.CreateClient(frame.ClientId, content.path, content.name)
	if err != nil {
		return err
	}
	return client.Clients.Add(&cl) // client.Clients.Add(content.name, &cl)
}

type registerCommandContent struct {
	name string
	path string
}

func parseData(data []byte) (registerCommandContent, error) {
	if len(data) < 4 {
		return registerCommandContent{}, errors.New("data must at least have delimiters")
	}

	nameLen := binary.BigEndian.Uint32(data[0:4])

	nameStartIdx := uint32(4)
	nameEndIdx := nameStartIdx + nameLen
	pathStartIdx := nameEndIdx

	if uint32(len(data)) < pathStartIdx {
		return registerCommandContent{}, errors.New("data to short")
	}

	nameRaw := data[nameStartIdx:nameEndIdx]
	pathRaw := data[pathStartIdx:]

	name := string(nameRaw)
	path := string(pathRaw)

	return registerCommandContent{
		name: name,
		path: path,
	}, nil
}
