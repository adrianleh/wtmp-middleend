package client

import (
	"errors"
	"fmt"
	"github.com/adrianleh/WTMP-middleend/messagequeue"
	"github.com/adrianleh/WTMP-middleend/types"
	"github.com/google/uuid"
	"sync"
)

type Client struct {
	id                    uuid.UUID
	socketPath            string
	name                  string
	acceptedTypes         []types.Type
	mqs                   map[types.Type]*messagequeue.MessageQueue
	superTypeMQCache      map[types.Type]*messagequeue.MessageQueue // Use to prevent having to walk type hierarchy
	superTypeMQCacheMutex *sync.RWMutex
	dataStructureMutex    *sync.Mutex
}

func CreateClient(id uuid.UUID, socketPath string, name string) Client {
	return Client{
		id:                 id,
		socketPath:         socketPath,
		name:               name,
		acceptedTypes:      make([]types.Type, 0),
		mqs:                map[types.Type]*messagequeue.MessageQueue{},
		superTypeMQCache:   map[types.Type]*messagequeue.MessageQueue{},
		dataStructureMutex: &sync.Mutex{},
	}
}

func (cl *Client) GetId() uuid.UUID               { return cl.id }
func (cl *Client) GetName() string                { return cl.name }
func (cl *Client) GetSocketPath() string          { return cl.socketPath }
func (cl *Client) GetAcceptedTypes() []types.Type { return cl.acceptedTypes }

type ClientMap struct {
	nameClientMap map[string]*Client // Contains UUID -> Socket mappings
	mutex         *sync.RWMutex
}

func CreateClientMap() ClientMap {
	return ClientMap{
		nameClientMap: map[string]*Client{},
		mutex:         &sync.RWMutex{},
	}
}

func (clients *ClientMap) Remove(name string) error {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	if clients.nameClientMap[name] == nil {
		return fmt.Errorf("client named \"%s\" does not exist", name)
	}
	clients.nameClientMap[name] = nil
	return nil
}

func (clients *ClientMap) Get(name string) *Client {
	clients.mutex.RLock()
	defer clients.mutex.RUnlock()
	return clients.nameClientMap[name]
}

func (clients *ClientMap) Add(name string, client *Client) error {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	if clients.nameClientMap[name] != nil {
		return fmt.Errorf("client named \"%s\" already exists", name)
	}
	clients.nameClientMap[name] = client
	return nil
}

var Clients = CreateClientMap()

func (cl *Client) Pop(typ types.Type) ([]byte, error) {
	if queue := cl.mqs[typ]; queue != nil {
		return queue.Pop()
	}
	return nil, fmt.Errorf("no queue found for type \"%s\"", typ.Name())
}

func (cl *Client) Empty(typ types.Type) bool {
	if queue := cl.mqs[typ]; queue != nil {
		return queue.Empty()
	}
	return false
}

func (cl *Client) Push(typ types.Type, data []byte) error {
	if queue := cl.getFromSuperTypeQueue(typ); queue != nil {
		return queue.Push(data)
	}
	superTypes := typ.GetSuperTypes()
	for _, superType := range superTypes {
		if queue := cl.mqs[superType]; queue != nil {
			cl.addToSuperTypeQueue(typ, queue)
			return queue.Push(data)
		}
	}
	return fmt.Errorf("no queue found for type \"%s\"", typ.Name())
}

func (cl *Client) RegisterType(typ types.Type) error {
	if cl.mqs[typ] != nil {
		return errors.New("type already registered")
	}
	cl.dataStructureMutex.Lock()
	cl.acceptedTypes = append(cl.acceptedTypes, typ)
	queue := messagequeue.CreateMessageQueue(typ.Size())
	cl.mqs[typ] = &queue
	cl.dataStructureMutex.Unlock()
	cl.invalidateSuperTypeCache()
	return nil
}

func (cl *Client) addToSuperTypeQueue(typ types.Type, mq *messagequeue.MessageQueue) {
	cl.superTypeMQCacheMutex.RLock() // We don't need a write lock here since overwriting is safe - as it would always be the same value
	defer cl.superTypeMQCacheMutex.RUnlock()
	cl.superTypeMQCache[typ] = mq
}

func (cl *Client) getFromSuperTypeQueue(typ types.Type) *messagequeue.MessageQueue {
	cl.superTypeMQCacheMutex.RLock()
	defer cl.superTypeMQCacheMutex.RUnlock()
	return cl.superTypeMQCache[typ]
}

func (cl *Client) invalidateSuperTypeCache() {
	cl.superTypeMQCacheMutex.Lock() // Once this is executed future reads are blocked until we unlock
	defer cl.superTypeMQCacheMutex.Unlock()
	cl.superTypeMQCache = map[types.Type]*messagequeue.MessageQueue{}
}
