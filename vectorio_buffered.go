/*
Copyright 2015 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
// Note: this is NOT concurrency safe.  Concurrent access should use the embedded Lock object (w.Lock.Lock() /
// w.Lock.Unlock()), or wrap this in a single goroutine that handles a channel of []byte.
type BufferedWritev struct {
	buf  []syscall.Iovec
	Lock *sync.Mutex
	fd   uintptr
}

// NewBufferedWritev makes a new BufferedWritev from a net.TCPConn, os.File, or file descriptor (FD).
func NewBufferedWritev(targetIn interface{}) (bw *BufferedWritev, err error) {
	switch target := targetIn.(type) {
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

// Write implements the io.Writer interface.
// Number of bytes written (nw) is usually 0 except for the times we flush the buffer, which will reflect the quantity of all bytes written in that writev() call
func (bw *BufferedWritev) Write(p []byte) (nw int, err error) {
	nw, err = bw.WriteIovec(syscall.Iovec{&p[0], uint64(len(p))})
	return
}

// WriteIovec takes a user-composed syscall.Iovec struct instead of []byte; functionally the same as Write
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

// Flush writes the contents of buffer to underlying file handle and resets buffer.
// This must be called at the end of writing before closing the underlying file, or data will be lost.
func (bw *BufferedWritev) Flush() (nw int, err error) {
	// TODO: if we're not going to use a lock, collapse Flush() and flush()
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
