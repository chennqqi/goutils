package utils

import (
	"io"
)

/*
MultiCopy like io.MultiWriters
only need an interface implement io.Write([]byte)
*/
func MultiCopy(src io.Reader, writers ...io.Writer) (written int64, err error) {
	return copyBuffer(nil, src, writers)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(buf []byte, src io.Reader, dsts []io.Writer) (written int64, err error) {
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			var nw int
			var ew error
			for _, dst := range dsts {
				nw, ew = dst.Write(buf[0:nr])
			}
			//only see last on
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
