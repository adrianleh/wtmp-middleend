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
	Id                    uuid.UUID
	SocketPath            string
	Name                  string
	AcceptedTypes         []types.Type
	MQs                   map[types.Type]*messagequeue.MessageQueue
	superTypeMQCache      map[types.Type]*messagequeue.MessageQueue // Use to prevent having to walk type hierarchy
	superTypeMQCacheMutex *sync.RWMutex
	dataStructureMutex    *sync.Mutex
}

func CreateClient(id uuid.UUID, socketPath string, name string) Client {
	return Client{
		Id:                 id,
		SocketPath:         socketPath,
		Name:               name,
		AcceptedTypes:      make([]types.Type, 0),
		MQs:                map[types.Type]*messagequeue.MessageQueue{},
		superTypeMQCache:   map[types.Type]*messagequeue.MessageQueue{},
		dataStructureMutex: &sync.Mutex{},
	}
}

type ClientMap map[string]*Client // Contains UUID -> Socket mappings

var Clients ClientMap = map[string]*Client{}

func (cl *Client) Pop(typ types.Type) ([]byte, error) {
	if queue := cl.MQs[typ]; queue != nil {
		return queue.Pop()
	}
	return nil, fmt.Errorf("no queue found for type \"%s\"", typ.Name())
}

func (cl *Client) Empty(typ types.Type) bool {
	if queue := cl.MQs[typ]; queue != nil {
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
		if queue := cl.MQs[superType]; queue != nil {
			cl.addToSuperTypeQueue(typ, queue)
			return queue.Push(data)
		}
	}
	return fmt.Errorf("no queue found for type \"%s\"", typ.Name())
}

func (cl *Client) RegisterType(typ types.Type) error {
	if cl.MQs[typ] != nil {
		return errors.New("type already registered")
	}
	cl.dataStructureMutex.Lock()
	cl.AcceptedTypes = append(cl.AcceptedTypes, typ)
	queue := messagequeue.CreateMessageQueue(typ.Size())
	cl.MQs[typ] = &queue
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
