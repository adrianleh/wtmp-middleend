package messagequeue

import (
	"sync"
	"testing"
)

func TestSeq(t *testing.T) {
	mq := CreateMessageQueue(1)
	b := []byte{1}
	if empty := mq.Empty(); !empty {
		t.Errorf("Should be empty init")
		return
	}
	if err := mq.Push(b); err != nil {
		t.Errorf("Failed to push, %v", err)
		return
	}
	if empty := mq.Empty(); empty {
		t.Errorf("Should not be empty")
		return
	}
	r, err := mq.Pop()
	if err != nil {
		t.Errorf("Failed to pop, %v", err)
		return
	}
	if len(r) != len(b) || r[0] != b[0] {
		t.Errorf("mismatch!")
		return
	}
	if empty := mq.Empty(); !empty {
		t.Errorf("Should be empty")
		return
	}
}

func TestPar(t *testing.T) {
	mq := CreateMessageQueue(4)
	b := []byte{1, 2, 3, 4}
	var wg sync.WaitGroup
	for i := 0; i < 2000000; i++ {
		wg.Add(1)
		go func() {
			_ = mq.Push(b)
			wg.Done()
		}()
	}
	wg.Wait()
	if len(mq.data) != 2000000 {
		t.Errorf("Didn't add enough items")
	}
	for i := 0; i < 2000000; i++ {
		wg.Add(1)
		go func() {
			_, _ = mq.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
	if len(mq.data) != 0 {
		t.Errorf("Didn't pop enough items")
	}
}
