package command

import (
	"errors"
	"strings"

	"github.com/adrianleh/WTMP-middleend/client"
)

type RegisterCommandHandler struct{}

func (RegisterCommandHandler) Handle(frame CommandFrame) error {
	content, err := parseData(frame.Data)
	if err != nil {
		return err
	}

	cl := client.CreateClient(frame.ClientId, content.path, content.name)

	return client.Clients.Add(&cl) // client.Clients.Add(content.name, &cl)
}

type registerCommandContent struct {
	name string
	path string
}

func parseData(data []byte) (registerCommandContent, error) {
	nullNameIdx := strings.Index(string(data), "\x00")
	if nullNameIdx < 0 {
		return registerCommandContent{}, errors.New("invalid format: no string found")
	}
	name := string(data[:nullNameIdx])
	if len(data) < nullNameIdx+1 {
		return registerCommandContent{}, errors.New("invalid format: data too short: no path")
	}
	rawPath := data[nullNameIdx+1:]
	nullPathIdx := strings.Index(string(rawPath), "\x00")
	path := string(rawPath[:nullPathIdx])

	// TODO: Length check entire frame

	return registerCommandContent{
		name: name,
		path: path,
	}, nil
}
