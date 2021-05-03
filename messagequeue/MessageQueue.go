package messagequeue

import (
	"errors"
	"sync"
)

type MessageQueue struct {
	elemSize uint64
	data     [][]byte
	lock     *sync.Mutex
}

func CreateMessageQueue(elemSize uint64) MessageQueue {
	return MessageQueue{
		elemSize: elemSize,
		data:     make([][]byte, 0),
		lock:     &sync.Mutex{},
	}
}

func (mq *MessageQueue) Empty() bool {
	return len(mq.data) == 0
}

func (mq *MessageQueue) Push(el []byte) error {
	if uint64(len(el)) != mq.elemSize {
		return errors.New("size mismatch")
	}
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.data = append(mq.data, el)
	return nil
}

func (mq *MessageQueue) Peek() ([]byte, error) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if len(mq.data) == 0 {
		return nil, errors.New("queue empty")
	}
	return mq.data[0], nil
}

func (mq *MessageQueue) Pop() ([]byte, error) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if len(mq.data) == 0 {
		return nil, errors.New("queue empty")
	}
	top := mq.data[0]
	mq.data = mq.data[1:]
	return top, nil
}
