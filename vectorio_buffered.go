package vectorio

import (
	"errors"
	"net"
	"os"
	"sync"
	"syscall"
)

// 1024 is the max size of an Iovec on Linux
const defaultBufSize = 1024

// just an alias so users don't have to import syscall to use this
//type Iovec syscall.Iovec

// BufferedWritev is similar to bufio.Writer.
// after all data has been written, the client should call Flush to make sure everything is written.
// Note: this is NOT concurrency safe.  Concurrent access should use the embedded Lock object (w.Lock.Lock() / w.Lock.Unlock()), or wrap this in a single goroutine that handles a channel of []byte.
type BufferedWritev struct {
	buf  []syscall.Iovec
	Lock *sync.Mutex
	fd   uintptr
}

// Make a new BufferedWritev from a net.TCPConn, os.File, or file descriptor (FD).
func NewBufferedWritev(target_in interface{}) (bw *BufferedWritev, err error) {
	switch target := target_in.(type) {
	case *net.TCPConn:
		var f *os.File
		f, err = target.File()
		if err != nil {
			return
		}
		bw, err = NewBufferedWritev(f)
	case *os.File:
		bw, err = NewBufferedWritev(uintptr(target.Fd()))
	case uintptr:
		// refactor: make buffer size user-specified?
		bw = &BufferedWritev{buf: make([]syscall.Iovec, 0, defaultBufSize), Lock: new(sync.Mutex), fd: target}
	default:
		err = errors.New("NewBufferedWritev called with invalid type")
	}
	return
}

// Implements the io.Writer interface.
// Number of bytes written (nw) is usually 0 except for the times we flush the buffer, which will reflect the quantity of all bytes written in that writev() call
func (bw *BufferedWritev) Write(p []byte) (nw int, err error) {
	nw, err = bw.WriteIovec(syscall.Iovec{&p[0], uint64(len(p))})
	return
}

func (bw *BufferedWritev) WriteIovec(iov syscall.Iovec) (nw int, err error) {
	//bw.lock.Lock()
	// normally append will reallocate a slice if it exceeds its cap, but that should not happen here because of our logic below
	bw.buf = append(bw.buf, iov)

	if len(bw.buf) == cap(bw.buf) {
		// maxed out the slice; write it and reset the slice
		nw, err = bw.flush()
	}

	//bw.lock.Unlock()
	return
}

// public interface; wraps flush() in a lock
// TODO: if we're not going to use a lock, collapse Flush() and flush()
func (bw *BufferedWritev) Flush() (nw int, err error) {
	//bw.Lock.Lock()
	nw, err = bw.flush()
	//bw.Lock.Unlock()
	return
}

// FUTURE: check to make sure the number of bytes written matches the sum of the iovec's
// Note: must be wrapped in a mutex to be concurrency safe
func (bw *BufferedWritev) flush() (nw int, err error) {
	nw, err = WritevRaw(bw.fd, bw.buf)
	bw.buf = bw.buf[:0]
	return
}
