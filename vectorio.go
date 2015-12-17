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
	nw, err = WritevRaw(uintptr(f.Fd()), iovec)
	return
}

// Call writev syscall, given a slice of syscall.Iovec to write
func WritevRaw(fd uintptr, iovec []syscall.Iovec) (nw int, err error) {
	nw_raw, _, errno := syscall.Syscall(syscall.SYS_WRITEV, fd, uintptr(unsafe.Pointer(&iovec[0])), uintptr(len(iovec)))
	nw = int(nw_raw)
	if errno != 0 {
		err = errors.New(fmt.Sprintf("writev failed with error: %d", errno))
	}
	return
}
