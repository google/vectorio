package syscalliovec

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// Call writev syscall, but first convert a [][]byte to []sycall.Iovec, return number of bytes written and an error
func Writev(f *os.File, in [][]byte) (nw int, err error) {
	iovec := make([]syscall.Iovec, len(in))
	for i, slice := range in {
		iovec[i] = syscall.Iovec{&slice[0], uint64(len(slice))}
	}
	nw, err = WritevRaw(f, iovec)
	return
}

// Call writev syscall, given a slice of syscall.Iovec to write
func WritevRaw(f *os.File, iovec []syscall.Iovec) (nw int, err error) {
	nw_raw, _, errno := syscall.Syscall(syscall.SYS_WRITEV, uintptr(f.Fd()), uintptr(unsafe.Pointer(&iovec[0])), uintptr(len(iovec)))
	nw = int(nw_raw)
	if errno != 0 {
		err = errors.New(fmt.Sprintf("writev failed with error: %d", errno))
	}
	return
}
