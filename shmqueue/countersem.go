package shmqueue

type CounterSem interface {
	Take() bool
	Give(bool)
	Destroy()
}

type CounterSemphore chan struct{}

func NewCounterSem(n int, init int) CounterSem {
	var c CounterSemphore
	c = make(chan struct{}, n)
	for i := 0; i < init; i++ {
		var a struct{}
		c <- a
	}
	return c
}

func (c CounterSemphore) Take() bool {
	_, ok := <-c
	return ok
}

func (c CounterSemphore) Give(force bool) {
	if force {
		var a struct{}
		c <- a
	} else {
		if len(c) == 0 {
			var a struct{}
			c <- a
		}
	}
}

func (c CounterSemphore) Destroy() {
	close(c)
}
