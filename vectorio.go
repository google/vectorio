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
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// Writev calls writev() syscall, but first convert a [][]byte to []sycall.Iovec, return number of bytes written and an error
func Writev(f *os.File, in [][]byte) (nw int, err error) {
	iovec := make([]syscall.Iovec, len(in))
	for i, slice := range in {
		iovec[i] = syscall.Iovec{&slice[0], uint64(len(slice))}
	}
	nw, err = WritevRaw(uintptr(f.Fd()), iovec)
	return
}

// WritevRaw calls writev() syscall like Writev, but expects a slice of syscall.Iovec
func WritevRaw(fd uintptr, iovec []syscall.Iovec) (nw int, err error) {
	nwRaw, _, errno := syscall.Syscall(syscall.SYS_WRITEV, fd, uintptr(unsafe.Pointer(&iovec[0])), uintptr(len(iovec)))
	nw = int(nwRaw)
	if errno != 0 {
		err = fmt.Errorf("writev failed with error: %d", errno)
	}
	return
}

// ReadvRaw calls the readv() syscall. The passed in iovec variable must be initialized
// to the expected number of Iovecs,  Len of each Iovec must be set,
// and the Base of each Iovec will be allocated if it is nil
func ReadvRaw(fd uintptr, iovec []syscall.Iovec) (num int, err error) {
	for i := range iovec {
		if iovec[i].Base == nil {
			base := make([]byte, iovec[i].Len)
			iovec[i].Base = &base[0]
		}
	}
	nrRaw, _, errno := syscall.Syscall(syscall.SYS_READV, fd, uintptr(unsafe.Pointer(&iovec[0])), uintptr(len(iovec)))
	num = int(nrRaw)
	if errno != 0 {
		err = fmt.Errorf("readv failed with error: %s", errno.Error())
	}
	return
}
