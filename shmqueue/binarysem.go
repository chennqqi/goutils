package shmqueue

type BinarySem interface {
	Take() bool
	Give(bool)
	Destroy()
}

type BinarySemphore chan struct{}

func NewBinarySem() BinarySem {
	var b BinarySemphore
	b = make(chan struct{}, 1)
	return b
}

func (b BinarySemphore) Take() bool {
	_, ok := <-b
	return ok
}

func (b BinarySemphore) Give(force bool) {
	if force {
		var a struct{}
		b <- a
	} else {
		if len(b) == 0 {
			var a struct{}
			b <- a
		}
	}
}

func (b BinarySemphore) Destroy() {
	close(b)
}
