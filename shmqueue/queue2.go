package shmqueue

import (
	"sync/atomic"
	"unsafe"
)

type ShmQueue2 interface {
	Push([]byte) error
	Pop([]byte) error
	Destroy()
}

type shmQueue2 struct {
	rIdx  uint64
	wIdx  uint64
	nDrop uint64
	max   uint64
	each  uint64

	buffer []byte

	sem CounterSem
}

func NewShmQueue2(segmentSize int, total int) (ShmQueue2, error) {
	var q shmQueue2
	q.max = uint64(total)
	q.buffer = make([]byte, (segmentSize+4)*total)
	q.sem = NewCounterSem(total, 0)

	q.each = uint64(segmentSize)
	return &q, nil
}

func (q *shmQueue2) Destroy() {
	q.sem.Destroy()
}

func (q *shmQueue2) Push(data []byte) error {
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

	//write data size
	valueOffset := int(offset * q.each)
	copy(q.buffer[valueOffset:], (*p1)[0:])
	copy(q.buffer[valueOffset+4:], data[:length])
	//fmt.Println("WRITE BUFFER LEN:", uLen, q.buffer[offset:offset+4])
	//fmt.Println("WRITE BUFFER OFFSET:", offset, q.buffer[valueOffset+4])

	//fmt.Println("w REMAIN:", remain)
	atomic.AddUint64(&q.wIdx, 1)
	sem.Give(true)
	return nil
}

func (q *shmQueue2) Pop(data []byte) error {
	sem := q.sem
	var err error
	for {
		if ok := sem.Take(); !ok {
			return ErrEOF
		}
		offset := q.rIdx % q.max

		var uLen uint32
		p := unsafe.Pointer(&uLen)
		p1 := (*[4]byte)(p)

		//read data size
		valueOffset := int(offset * q.each)
		copy((*p1)[0:], q.buffer[valueOffset:])

		capLen := uint64(cap(data))
		if capLen > q.each {
			capLen = q.each
		} else if capLen < uint64(uLen) {
			err = ErrTruncate
		}

		copy(data, q.buffer[valueOffset+4:valueOffset+4+int(capLen)])
		break
	}

	atomic.AddUint64(&q.rIdx, 1)
	return err
}

func (q *shmQueue2) PopNoneblock(data []byte) error {
	remain := q.wIdx - q.rIdx
	if remain == 0 {
		return ErrBlockR
	}

	var err error
	offset := q.rIdx % q.max
	valueOfset := int(offset * q.each)

	var uLen uint32
	p := unsafe.Pointer(&uLen)
	p1 := (*[4]byte)(p)

	//read data size
	copy((*p1)[0:], q.buffer[valueOfset:])

	capLen := uint64(cap(data))
	if capLen > q.each {
		capLen = q.each
	} else if capLen < uint64(uLen) {
		err = ErrTruncate
	}
	copy(data, q.buffer[valueOfset+4:valueOfset+4+int(capLen)])
	return err
}
