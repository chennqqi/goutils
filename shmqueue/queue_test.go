package shmqueue

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	MAX_TEST = 100000
)

func TestQueue(t *testing.T) {
	q, err := New(1024, 256)
	assert.Nil(t, err)
	t.Parallel()

	quit := make(chan struct{})
	go func() {
		var buf [16]byte
		var zero [16]byte
		var count int
		for {
			copy(buf[:], zero[:])
			err := q.Pop(buf[:])
			if err == ErrEOF {
				fmt.Println("QUIT")
				close(quit)
				return
			}
			if err != nil {
				t.Fatal("Expect Pop Error:", err)
			}
			value := fmt.Sprintf("%d", count)
			//			fmt.Println("RX:", string(buf[:]), value)
			for i := 0; i < len(value); i++ {
				if value[i] != buf[i] {
					t.Error("count", count, "not equal")
				}
				//				assert.Equal(t, value[i], buf[i])
			}
			count++
		}
	}()

	var block int
	for i := 0; i < MAX_TEST; i++ {
		value := fmt.Sprintf("%d", i)
		for {
			err := q.Push([]byte(value))
			if err == ErrBlockW {
				runtime.Gosched()
				block++
			} else {
				break
			}
		}
	}
	t.Log("block:", block)
	q.Destroy()
	<-quit
}

func BenchmarkQueue(b *testing.B) {
	q, err := New(1024, 256)
	assert.Nil(b, err)

	quit := make(chan struct{})
	go func() {
		var buf [4]byte
		var zero [4]byte
		var count int
		for {
			copy(buf[:], zero[:])
			err := q.Pop(buf[:])
			if err == ErrEOF {
				close(quit)
				return
			}
			count++
		}
	}()

	var block int
	for i := 0; i < b.N; i++ {
		value := fmt.Sprintf("%d", i)
		for {
			err := q.Push([]byte(value))
			if err == ErrBlockW {
				runtime.Gosched()

				block++
			} else {
				break
			}
		}
	}
	b.Log("block:", block)
	q.Destroy()

	<-quit
}

func TestQueue2(t *testing.T) {
	return
	t.Parallel()
	q, err := NewShmQueue2(1024, 256)
	assert.Nil(t, err)

	quit := make(chan struct{})
	go func() {
		var buf [16]byte
		var zero [16]byte
		var count int
		for {
			copy(buf[:], zero[:])
			err := q.Pop(buf[:])
			if err == ErrEOF {
				close(quit)
				return
			}
			if err != nil {
				t.Fatal("Expect Pop Error:", err)
			}
			value := fmt.Sprintf("%d", count)
			for i := 0; i < len(value); i++ {
				if value[i] != buf[i] {
					t.Fatal("count", count, "not equal")
				}
			}
			count++
		}
	}()

	var block int
	for i := 0; i < MAX_TEST; i++ {
		value := fmt.Sprintf("%d", i)
		for {
			err := q.Push([]byte(value))
			if err == ErrBlockW {
				runtime.Gosched()
				block++
			} else {
				//		fmt.Println("TX", value)
				break
			}
		}
	}
	fmt.Println("block:", block)
	q.Destroy()
	<-quit
}

func BenchmarkQueue2(b *testing.B) {
	return
	q, err := NewShmQueue2(1024, 256)
	assert.Nil(b, err)

	quit := make(chan struct{})
	go func() {
		var buf [16]byte
		var zero [16]byte
		var count int
		for {
			copy(buf[:], zero[:])
			err := q.Pop(buf[:])
			if err == ErrEOF {
				close(quit)
				return
			}
			if err != nil {
				b.Fatal("Expect Pop Error:", err)
			}
			count++
		}
	}()

	var block int
	for i := 0; i < b.N; i++ {
		value := fmt.Sprintf("%d", i)
		for {
			err := q.Push([]byte(value))
			if err == ErrBlockW {
				runtime.Gosched()

				block++
			} else {
				break
			}
		}
	}
	b.Log("block:", block)
	q.Destroy()

	<-quit
}

//func BenchmarkLoopsParallel(b *testing.B) {
//    b.RunParallel(func(pb *testing.PB) {
//        var test ForTest
//        ptr := &test
//        for pb.Next() {
//            ptr.Loops()
//        }
//    }
//}
