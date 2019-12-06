package ringstore

import (
	"fmt"
	"math"
	"sync"
)

type MemoryRingStore struct {
	locker     sync.RWMutex
	items      []interface{}
	length     uint64
	start, end uint64
}

func NewMemoryRingStore(capacity uint64) *MemoryRingStore {
	r := &MemoryRingStore{
		items:  make([]interface{}, capacity, capacity),
		length: capacity,
		start:  0,
		end:    0,
	}
	return r
}

func (mrb *MemoryRingStore) Append(data interface{}) (err error) {
	mrb.locker.Lock()
	defer mrb.locker.Unlock()
	mrb.items[mrb.end%mrb.length] = data
	mrb.end += 1

	// not overflow
	if mrb.end > mrb.start {
		if mrb.end-mrb.start >= mrb.length {
			mrb.start += 1
		}
	} else if mrb.end < mrb.start && mrb.end+(math.MaxUint64-mrb.start+1) >= mrb.length {
		mrb.start += 1
	} else {
		panic(fmt.Sprintf("unexpected state: (start, end, length) = (%d, %d, %d)", mrb.start, mrb.end, mrb.length))
	}
	return
}

func (mrb *MemoryRingStore) Get(c *Cursor) (data interface{}, state RingStoreGetState, err error) {
	mrb.locker.RLock()
	defer mrb.locker.RUnlock()

	// no data in ring buffer
	if mrb.start == mrb.end {
		return nil, Empty, nil
	}

	// first read
	if !c.tmpInitial {
		c.tmpInitial = true
		c.tmpNext = mrb.start
	}
	if c.tmpNext == mrb.end {
		return nil, Beyond, nil
	}
	if mrb.start < mrb.end {
		if c.tmpNext < mrb.start {
			c.tmpNext = mrb.start + 1
			return mrb.items[mrb.start%mrb.length], Behind, nil
		}
	} else {
		if c.tmpNext < mrb.start && c.tmpNext > mrb.end {
			c.tmpNext = mrb.start + 1
			return mrb.items[mrb.start%mrb.length], Behind, nil
		}
	}
	// normal read and move tmp cursor
	curr := c.tmpNext
	c.tmpNext++
	return mrb.items[curr%mrb.length], Match, nil
}

func minUint64(a, b uint64) uint64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func (mrb *MemoryRingStore) Close() error {
	return nil
}
