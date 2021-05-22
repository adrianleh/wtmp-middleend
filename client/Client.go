package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/adrianleh/WTMP-middleend/messagequeue"
	"github.com/adrianleh/WTMP-middleend/types"
	"github.com/google/uuid"
	"io"
	"net"
	"sync"
)

type Client struct {
	id                    uuid.UUID
	socketPath            string
	name                  string
	acceptedTypes         []types.Type
	mqs                   map[string]*messagequeue.MessageQueue
	superTypeCache        map[string]*types.Type
	superTypeCacheMutex   *sync.RWMutex
	dataStructureMutex    *sync.Mutex
	sock                  net.Conn
	inOrderExecutionMutex *sync.Mutex
}

func CreateClient(id uuid.UUID, socketPath string, name string) (Client, error) {
	sock, err := net.Dial("unix", socketPath)
	if err != nil {
		return Client{}, err
	}
	return Client{
		id:                    id,
		socketPath:            socketPath,
		name:                  name,
		acceptedTypes:         make([]types.Type, 0),
		mqs:                   map[string]*messagequeue.MessageQueue{},
		superTypeCache:        map[string]*types.Type{},
		dataStructureMutex:    &sync.Mutex{},
		superTypeCacheMutex:   &sync.RWMutex{},
		sock:                  sock,
		inOrderExecutionMutex: &sync.Mutex{},
	}, nil
}

func (cl *Client) SendToClient(data []byte) error {
	_, err := io.Copy(cl.sock, bytes.NewReader(data))
	return err
}

func (cl *Client) GetCommandMutex() *sync.Mutex {
	return cl.inOrderExecutionMutex
}

func (cl *Client) GetId() uuid.UUID               { return cl.id }
func (cl *Client) GetName() string                { return cl.name }
func (cl *Client) GetSocketPath() string          { return cl.socketPath }
func (cl *Client) GetAcceptedTypes() []types.Type { return cl.acceptedTypes }

type ClientMap struct {
	uuidClientMap map[uuid.UUID]*Client
	nameClientMap map[string]*Client
	mutex         *sync.RWMutex
}

func CreateClientMap() ClientMap {
	return ClientMap{
		uuidClientMap: map[uuid.UUID]*Client{},
		nameClientMap: map[string]*Client{},
		mutex:         &sync.RWMutex{},
	}
}

func (clients *ClientMap) Remove(id uuid.UUID) error {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	if clients.uuidClientMap[id] == nil {
		return fmt.Errorf("client with id \"%s\" does not exist", id.String())
	}
	name := clients.uuidClientMap[id].GetName()
	clients.nameClientMap[name] = nil
	clients.uuidClientMap[id] = nil
	return nil
}

func (clients *ClientMap) GetByName(name string) *Client {
	clients.mutex.RLock()
	defer clients.mutex.RUnlock()
	return clients.nameClientMap[name]
}

func (clients *ClientMap) GetById(id uuid.UUID) *Client {
	clients.mutex.RLock()
	defer clients.mutex.RUnlock()
	return clients.uuidClientMap[id]
}

func (clients *ClientMap) Add(client *Client) error {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	name := client.GetName()
	if clients.nameClientMap[name] != nil {
		return fmt.Errorf("client named \"%s\" already exists", name)
	}
	clients.nameClientMap[name] = client
	clients.uuidClientMap[client.GetId()] = client
	return nil
}

var Clients = CreateClientMap()

func (cl *Client) Pop(typ types.Type) ([]byte, error) {
	if queue := cl.mqs[typ.Name()]; queue != nil {
		return queue.Pop()
	}
	return nil, fmt.Errorf("no queue found for type \"%s\"", typ.Name())
}

func (cl *Client) Empty(typ types.Type) (bool, error) {
	if queue := cl.mqs[typ.Name()]; queue != nil {
		return queue.Empty(), nil
	}
	return false, errors.New("no queue exists for type")
}

func (cl *Client) PushToSuperType(typ types.Type, superType types.Type, data []byte) error {
	trimmedData, err := types.Trim(typ, superType, data)
	if err != nil {
		return err
	}
	queue := cl.mqs[superType.Name()]
	return queue.Push(trimmedData)
}

func (cl *Client) Push(typ types.Type, data []byte) error {
	if superType := cl.getFromSuperTypeCache(typ); superType != nil {
		return cl.PushToSuperType(typ, *superType, data)
	}
	superTypes := typ.GetSuperTypes()
	for _, superType := range superTypes {
		if queue := cl.mqs[superType.Name()]; queue != nil {
			cl.addToSuperTypeCache(typ, &superType)
			return cl.PushToSuperType(typ, superType, data)
		}
	}
	return fmt.Errorf("no queue found for type \"%s\"", typ.Name())
}

func (cl *Client) RegisterType(typ types.Type) error {
	if cl.mqs[typ.Name()] != nil {
		return errors.New("type already registered")
	}
	cl.dataStructureMutex.Lock()
	cl.acceptedTypes = append(cl.acceptedTypes, typ)
	queue := messagequeue.CreateMessageQueue(typ.Size())
	cl.mqs[typ.Name()] = &queue
	cl.dataStructureMutex.Unlock()
	cl.invalidateSuperTypeCache()
	return nil
}

func (cl *Client) addToSuperTypeCache(typ types.Type, super *types.Type) {
	cl.superTypeCacheMutex.RLock() // We don't need a write lock here since overwriting is safe - as it would always be the same value
	defer cl.superTypeCacheMutex.RUnlock()
	cl.superTypeCache[typ.Name()] = super
}

func (cl *Client) getFromSuperTypeCache(typ types.Type) *types.Type {
	cl.superTypeCacheMutex.RLock()
	defer cl.superTypeCacheMutex.RUnlock()
	return cl.superTypeCache[typ.Name()]
}

func (cl *Client) invalidateSuperTypeCache() {
	cl.superTypeCacheMutex.Lock() // Once this is executed future reads are blocked until we unlock
	defer cl.superTypeCacheMutex.Unlock()
	cl.superTypeCache = map[string]*types.Type{}
}
