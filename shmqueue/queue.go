package shmqueue

import (
	"errors"
	//	"fmt"
	"sync/atomic"
	"unsafe"
)

var (
	ErrBlockR   = errors.New("BlockR")
	ErrBlockW   = errors.New("BlockW")
	ErrEOF      = errors.New("EOF")
	ErrTruncate = errors.New("truncate")
)

type ShmQueue interface {
	Push([]byte) error
	Pop([]byte) error
	Destroy()
}

type shmQueue struct {
	rIdx  uint64
	wIdx  uint64
	nDrop uint64
	max   uint64
	each  uint64

	buffer []byte

	sem BinarySem
}

func New(segmentSize int, total int) (ShmQueue, error) {
	var q shmQueue
	q.max = uint64(total)
	q.buffer = make([]byte, segmentSize*total)
	q.sem = NewBinarySem()
	q.each = uint64(segmentSize)
	return &q, nil
}

func (q *shmQueue) Destroy() {
	q.sem.Destroy()
}

func (q *shmQueue) Push(data []byte) error {
	remain := q.wIdx - q.rIdx
	sem := q.sem

	if remain == q.max {
		q.nDrop++
		return ErrBlockW
	}
	// do copy

	offset := q.wIdx % q.max

	length := len(data)
	if uint64(length) > q.each {
		length = int(q.each)
	}

	var uLen = uint32(length)
	p := unsafe.Pointer(&uLen)
	p1 := (*[4]byte)(p)

	//fmt.Println("length:", length, uLen, (*p1)[0:])
	//write data size
	valueOffset := int(offset * q.each)
	copy(q.buffer[valueOffset:], (*p1)[0:])
	copy(q.buffer[valueOffset+4:], data[:length])
	//fmt.Println("WRITE BUFFER LEN:", uLen, q.buffer[offset:offset+4])
	//fmt.Println("WRITE BUFFER OFFSET:", offset, q.buffer[valueOffset+4])

	//fmt.Println("w REMAIN:", remain)
	atomic.AddUint64(&q.wIdx, 1)
	if remain == 1 {
		//fmt.Println("semGive")
		sem.Give(false)
	}
	return nil
}

func (q *shmQueue) Pop(data []byte) error {
	sem := q.sem
	var err error
	for {
		remain := q.wIdx - q.rIdx
		//fmt.Println("r REMAIN:", remain)
		if remain == 0 {
			//fmt.Println("pop BLOCK")
			if ok := sem.Take(); !ok {
				//fmt.Println("EOF")
				return ErrEOF
			}
			//fmt.Println("pop RESUME")
		} else {
			offset := q.rIdx % q.max

			var uLen uint32
			p := unsafe.Pointer(&uLen)
			p1 := (*[4]byte)(p)

			//read data size
			valueOffset := int(offset * q.each)
			copy((*p1)[0:], q.buffer[valueOffset:])

			capLen := uint64(cap(data))

			//			fmt.Println("OFFSET:", offset, q.buffer[valueOffset])
			//			fmt.Println("OFFSET DAT LEN:", uLen)

			if capLen > q.each {
				capLen = q.each
			} else if capLen < uint64(uLen) {
				err = ErrTruncate
			}

			copy(data, q.buffer[valueOffset+4:valueOffset+4+int(capLen)])
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
