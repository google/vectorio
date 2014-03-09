package syscalliovec

import (
	"net"
	"os"
	"sync"
	"syscall"
)

// 1024 is the max size of an Iovec on Linux
const defaultBufSize = 1024

// just an alias so users don't have to import syscall to use this
//type Iovec syscall.Iovec

// BufferedWritev is similar to bufio.Writer
// after all data has been written, the client should call Flush to make sure everything is written
type BufferedWritev struct {
	buf  []syscall.Iovec
	lock *sync.Mutex
	fd   uintptr
}

// Make a new BufferedWritev from a net.TCPConn (null if err)
func NewBufferedWritevTCPConn(tcp *net.TCPConn) (bw *BufferedWritev, err error) {
	f, err := tcp.File()
	if err != nil {
		return
	}
	bw = NewBufferedWritevFile(f)
	return
}

// Make a new BufferedWritev from a os.File
func NewBufferedWritevFile(f *os.File) *BufferedWritev {
	return NewBufferedWritevFD(uintptr(f.Fd()))
}

// Make a new BufferedWritev from a file descriptor (FD)
func NewBufferedWritevFD(fd uintptr) *BufferedWritev {
	// refactor: make buffer size user-specified?
	return &BufferedWritev{buf: make([]syscall.Iovec, 0, defaultBufSize), lock: new(sync.Mutex), fd: fd}
}

// this implements the io.Writer interface
// except we always return 0 bytes, which should be ignored
func (bw *BufferedWritev) Write(p []byte) (n int, err error) {
	err = bw.WriteIovec(syscall.Iovec{&p[0], uint64(len(p))})
	return
}

func (bw *BufferedWritev) WriteIovec(iov syscall.Iovec) (err error) {
	bw.lock.Lock()
	// normally append will reallocate a slice if it exceeds its cap, but that should not happen here because of our logic below
	bw.buf = append(bw.buf, iov)

	if len(bw.buf) == cap(bw.buf) {
		// maxed out the slice; write it and reset the slice
		err = bw.flush()
	}

	bw.lock.Unlock()
	return
}

// public interface; wraps flush() in a lock
func (bw *BufferedWritev) Flush() (err error) {
	bw.lock.Lock()
	err = bw.flush()
	bw.lock.Unlock()
	return
}

// FUTURE: check to make sure the number of bytes written matches the sum of the iovec's
// Note: must be wrapped in a mutex to be concurrency safe
func (bw *BufferedWritev) flush() (err error) {
	_, err = WritevRaw(bw.fd, bw.buf)
	bw.buf = bw.buf[:0]
	return
}
