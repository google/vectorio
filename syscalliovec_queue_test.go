package syscalliovec

import (
	"bytes"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func TestBufferedWritev(t *testing.T) {
	f, err := os.Create("testwritevraw")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := []byte("foobarbaz")
	data_desired := []byte("foobazbar")

	bw := NewBufferedWritevFile(f)

	bw.WriteIovec(syscall.Iovec{&data[0], 3})
	bw.WriteIovec(syscall.Iovec{&data[6], 3})
	bw.Write(data[3:6])

	err = bw.Flush()
	if err != nil {
		t.Errorf("Flush threw error: %s", err)
	}

	stats, err := f.Stat()
	if err != nil {
		t.Errorf("Stat threw error: %s", err)
	}
	nw := stats.Size()
	if nw != int64(len(data)) {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), nw)
	}

	f.Seek(0, 0)
	fromdisk, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read file back, %s", err)
	}
	if bytes.Compare(fromdisk, data_desired) != 0 {
		t.Errorf("contents of file don't match input, %s != %s or %d != %d", fromdisk, data_desired, len(fromdisk), len(data_desired))
	}
	os.Remove("testwritevraw")
}
