package client

import "github.com/google/uuid"

type Client struct {
	Id         uuid.UUID
	SocketPath string
	Name       string
}

type ClientMap map[string]*Client // Contains UUID -> Socket mappings

var Clients ClientMap = map[string]*Client{}
