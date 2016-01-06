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

package vectorio_test

import (
	"fmt"
	"io/ioutil"
	"syscall"

	"github.com/google/vectorio"
)

func ExampleBufferedVectorio() {
	// Create a temp file to use
	f, err := ioutil.TempFile("", "vectorio")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data1 := []byte("foobarbaz")
	data2 := []byte("foobazbar")

	w, err := vectorio.NewBufferedWritev(f)
	nw, err := w.Write(data2)
	nw, err = w.WriteIovec(syscall.Iovec{&data1[0], 9})
	nw, err = w.Flush()

	if err != nil {
		fmt.Println("Flush threw error: %s", err)
	}
	if nw == 9*2 {
		fmt.Println("Wrote", nw, "bytes to file!")
	} else {
		fmt.Println("did not write 9 * 2 bytes, wrote ", nw)
	}

	// Output:
	// Wrote 18 bytes to file!
}

func ExampleBufferedVectorioUnsafe() {
	// Create a temp file to use
	f, _ := ioutil.TempFile("", "vectorio")
	defer f.Close()

	data1 := []byte("foobarbaz")
	data2 := []byte("foobazbar")

	w, _ := vectorio.NewBufferedWritev(f)
	w.Write(data2)
	w.WriteIovec(syscall.Iovec{&data1[0], 9})
	nw, _ := w.Flush()
	fmt.Println("Wrote", nw, "bytes to file!")

	// Output:
	// Wrote 18 bytes to file!
}
