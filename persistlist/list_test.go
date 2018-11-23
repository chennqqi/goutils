package persistlist

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushPop(t *testing.T) {
	list, err := NewNodbList("test", "list")
	assert.Nil(t, err)

	var value int
	for i := 0; i < 100; i++ {
		value = i
		acount, err := list.Push(value)
		assert.Equal(t, acount, int64(i+1))
		assert.Nil(t, err)
	}

	for i := 0; i < 100; i++ {
		err = list.Pop(&value)
		assert.Nil(t, err)
		assert.Equal(t, value, i)
	}
	err = list.Pop(&value)
	assert.Equal(t, err, ErrNil)

	list.Close()
	err = os.RemoveAll("test")
	assert.Nil(t, err)
}

func BenchmarkAPush(b *testing.B) {
	list, err := NewNodbList("test", "list")
	assert.Nil(b, err)
	assert.NotNil(b, list)
	defer list.Close()

	var value int
	for i := 0; i < b.N; i++ {
		value = i
		_, err := list.Push(value)
		assert.Nil(b, err)
	}
}

func BenchmarkBPop(b *testing.B) {
	list, err := NewNodbList("test", "list")
	assert.Nil(b, err)

	var value int
	for i := 0; i < b.N; i++ {
		value = i
		_, err := list.Push(value)
		assert.Nil(b, err)
	}
	list.Close()
	os.RemoveAll("test")
}
