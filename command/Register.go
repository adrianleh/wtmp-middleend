package command

import (
	"errors"
	"fmt"
	"github.com/adrianleh/WTMP-middleend/client"
	"strings"
)

type RegisterCommandHandler struct{}

func (RegisterCommandHandler) Handle(frame CommandFrame) error {
	content, err := parseData(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Client{
		Id:         frame.ClientId,
		SocketPath: content.path,
		Name:       content.name,
	}

	if client.Clients[content.name] != nil {
		return fmt.Errorf("client \"%s\" already registered", content.name)
	}

	client.Clients[content.name] = &cl
	return nil
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
