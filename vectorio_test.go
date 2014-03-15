package vectorio

import (
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWritevRaw(t *testing.T) {
	f, err := os.Create("testwritevraw")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := []byte("foobarbaz")
	data_desired := []byte("foobazbar")
	iovec := []syscall.Iovec{syscall.Iovec{&data[0], 3}, syscall.Iovec{&data[6], 3}, syscall.Iovec{&data[3], 3}}

	nw, err := WritevRaw(uintptr(f.Fd()), iovec)
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
	if bytes.Compare(fromdisk, data_desired) != 0 {
		t.Errorf("contents of file don't match input, %s != %s or %d != %d", fromdisk, data_desired, len(fromdisk), len(data_desired))
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

func TestWritevSocket(t *testing.T) {
	go func() {
		ln, err := net.Listen("tcp", "127.0.0.1:9999")
		if err != nil {
			t.Errorf("could not listen on 127.0.0.1:9999: %s", err)
		}

		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("could not accept conn: %s", err)
		}
		defer conn.Close()

		buf := make([]byte, 9)
		nr, err := conn.Read(buf)
		if nr != len(buf) {
			t.Errorf("read was wrong length: %d != %d", nr, len(buf))
		}

		good := []byte("foobarbaz")
		if bytes.Compare(buf, good) != 0 {
			t.Errorf("read got wrong data: %s != %s", buf, good)
		} else {
			t.Logf("read got correct data: %s == %s", buf, good)
		}
	}()

	time.Sleep(1 * time.Second)
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:9999")
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		t.Errorf("error connecting to 127.0.0.1:9999: %s", err)
	} else {
		t.Logf("connected to server")
	}
	defer conn.Close()
	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}

	f, err := conn.File()
	if err != nil {
		t.Errorf("could not get file handle for TCP client: %s", err)
	}
	defer f.Close()
	nw, err := Writev(f, data)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != 9 {
		t.Errorf("Length %d of input does not match %d written bytes", len(data), nw)
	} else {
		t.Logf("right number of bytes written")
	}
	time.Sleep(1 * time.Second)
}
