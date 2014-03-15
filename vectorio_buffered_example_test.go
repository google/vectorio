package vectorio_test

import (
	"fmt"
	"github.com/tedb/vectorio"
	"io/ioutil"
	"syscall"
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
