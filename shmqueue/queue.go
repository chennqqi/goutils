package shmqueue

import (
	"errors"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	ErrBlockR   = errors.New("BlockR")
	ErrBlockW   = errors.New("BlockW")
	ErrTruncate = errors.New("truncate")
)

type ShmQueue interface {
	Push([]byte) error
	Pop([]byte) error
}

type shmQueue struct {
	rIdx  uint64
	wIdx  uint64
	nDrop uint64
	max   uint64
	each  uint64

	buffer []byte

	cond  *sync.Cond
	mutex sync.Mutex
}

func New(segmentSize int, total int) (ShmQueue, error) {
	var q shmQueue
	q.cond = sync.NewCond(&q.mutex)
	return &q, nil
}

func (q *shmQueue) Push(data []byte) error {
	remain := q.wIdx - q.rIdx
	cond := q.cond

	if remain == q.max {
		q.nDrop++
		return ErrBlockW
	}
	// do copy
	atomic.AddUint64(&q.wIdx, 1)

	offset := q.wIdx % q.max

	length := len(data)
	if uint64(length) > q.each {
		length = int(q.each)
	}

	var uLen = uint32(length)
	p := unsafe.Pointer(&uLen)
	p1 := (*[4]byte)(p)

	//write data size
	copy(q.buffer[offset:], (*p1)[0:])
	copy(q.buffer[offset+4:], data[:length])

	if remain == 0 {
		cond.Signal()
	}
	return nil
}

func (q *shmQueue) Pop(data []byte) error {
	cond := q.cond
	var err error
	for {
		remain := q.wIdx - q.rIdx
		if remain == 0 {
			cond.Wait()
		} else {
			offset := q.rIdx % q.max

			var uLen uint32
			p := unsafe.Pointer(&uLen)
			p1 := (*[4]byte)(p)

			//read data size
			copy((*p1)[0:], q.buffer[offset:])

			capLen := uint64(cap(data))
			if capLen > q.each {
				capLen = q.each
			} else if capLen < uint64(uLen) {
				err = ErrTruncate
			}

			copy(data, q.buffer[offset+4:offset+4+capLen])
			break
		}
	}

	atomic.AddUint64(&q.rIdx, 1)
	return err
}

func (q *shmQueue) PopNoneblock(data []byte) error {
	remain := q.wIdx - q.rIdx
	if remain == 0 {
		return ErrBlockR
	}

	var err error
	offset := q.rIdx % q.max

	var uLen uint32
	p := unsafe.Pointer(&uLen)
	p1 := (*[4]byte)(p)

	//read data size
	copy((*p1)[0:], q.buffer[offset:])

	capLen := uint64(cap(data))
	if capLen > q.each {
		capLen = q.each
	} else if capLen < uint64(uLen) {
		err = ErrTruncate
	}
	copy(data, q.buffer[offset+4:offset+4+capLen])
	return err
}
