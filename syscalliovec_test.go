package syscalliovec

import (
	"bytes"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func TestWritevRaw(t *testing.T) {
	f, err := os.Create("testwritevraw")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := []byte("foobarbaz")
	iovec := []syscall.Iovec{syscall.Iovec{&data[0], 3}, syscall.Iovec{&data[3], 3}, syscall.Iovec{&data[6], 3}}

	nw, err := WritevRaw(f, iovec)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != len(data) {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), nw)
	}

	fromdisk, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read file back, %s", err)
	}
	if bytes.Compare(fromdisk, data) != 0 {
		t.Errorf("contents of file don't match input, %s != %s, %d != %d", fromdisk, data, len(fromdisk), len(data))
	}
	os.Remove("testwritevraw")
}

func TestWritev(t *testing.T) {
	f, err := os.Create("testwritev")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}

	nw, err := Writev(f, data)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != 9 {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), nw)
	}

	fromdisk, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read file back, %s", err)
	}
	should := []byte("foobarbaz")
	if bytes.Compare(fromdisk, should) != 0 {
		t.Errorf("contents of file don't match input, %s != %s", fromdisk, should)
	}
	os.Remove("testwritev")
}
