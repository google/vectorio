package vectorio

import (
	"bytes"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func TestBufferedWritev(t *testing.T) {
	f, err := os.Create("testbufferedwritev")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := []byte("foobarbaz")
	data_desired := []byte("foobazbar")

	bw, _ := NewBufferedWritev(f)

	bw.WriteIovec(syscall.Iovec{&data[0], 3})
	bw.WriteIovec(syscall.Iovec{&data[6], 3})
	bw.Write(data[3:6])

	nw, err := bw.Flush()
	if err != nil {
		t.Errorf("Flush threw error: %s", err)
	}
	if nw != 3*3 {
		t.Errorf("Flush wrote wrong number of bytes, expected 3 * 3 but got", nw)
	}

	stats, err := f.Stat()
	if err != nil {
		t.Errorf("Stat threw error: %s", err)
	}
	filesize := stats.Size()
	if filesize != int64(len(data)) {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), filesize)
	}

	f.Seek(0, 0)
	fromdisk, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read file back, %s", err)
	}
	if bytes.Compare(fromdisk, data_desired) != 0 {
		t.Errorf("contents of file don't match input, %s != %s or %d != %d", fromdisk, data_desired, len(fromdisk), len(data_desired))
	}
	os.Remove("testbufferedwritev")
}

func TestBufferedWritevHuge(t *testing.T) {
	f, err := os.Create("testbufferedwritevhuge")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := []byte("foobarba")
	final := []byte("attheend")

	bw, _ := NewBufferedWritev(f)

	// write a little more than our buffer size
	for i := 0; i < 1024; i++ {
		bw.Write(data)
	}
	bw.Write(final)

	_, err = bw.Flush()
	if err != nil {
		t.Errorf("Flush threw error: %s", err)
	}

	stats, err := f.Stat()
	if err != nil {
		t.Errorf("Stat threw error: %s", err)
	}
	filesize := stats.Size()
	if filesize != int64(1024*len(data)+len(final)) {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), filesize)
	}

	// Maybe make this validate file contents later??
	f.Seek(0, 0)
	fromdisk, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read file back, %s", err)
	}
	compared := fromdisk[1024*8:]
	if bytes.Compare(compared, final) != 0 {
		t.Errorf("contents of file don't match input, %s != %s or %d != %d", compared, final, len(compared), len(final))
	}
	os.Remove("testbufferedwritevhuge")
}
