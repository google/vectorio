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
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/google/vectorio"
)

func TestWritevRaw(t *testing.T) {
	f, err := os.Create("testwritevraw")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	iovec := []syscall.Iovec{
		syscall.Iovec{Base:(*byte)(unsafe.Pointer(&data[0][0])), Len:uint64(len(data[0]))},
		syscall.Iovec{Base:(*byte)(unsafe.Pointer(&data[1][0])), Len:uint64(len(data[1]))},
		syscall.Iovec{Base:(*byte)(unsafe.Pointer(&data[2][0])), Len:uint64(len(data[2]))}}

	iov := make([]syscall.Iovec, len(data))
	expected := 0
	for i := range data {
		expected += len(data[i])
		iov[i].Len = uint64(len(data[i]))
	}

	nw, err := vectorio.WritevRaw(uintptr(f.Fd()), iovec)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != expected {
		t.Errorf("Length %d of input does not match %d written bytes", expected, nw)
	}

	sz, err := vectorio.ReadvRaw(uintptr(f.Fd()), iov)
	if err != nil {
		t.Errorf("read returned error %s", err.Error())
	}

	if sz != expected {
		t.Errorf("read was wrong length: %d != %d", sz, expected)
	}

	for i := range data {
		b := *(*[]byte)(unsafe.Pointer(&iov[i].Base))
		if bytes.Compare(data[i], b) != 0 {
			t.Errorf("read got wrong data: %s != %s", data[i], b)
		} else {
			t.Logf("read got correct data: %s == %s", data[i], b)
		}
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

	nw, err := vectorio.Writev(f, data)
	f.Seek(0, 0)
	if err != nil {
		t.Errorf("WritevRaw threw error: %s", err)
	}

	if nw != 9 {
		t.Errorf("Length %d of input does not match %d written bytes", 9, nw)
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
	data := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}

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

		iov := make([]syscall.Iovec, len(data))

		expected := 0
		for i := range data {
			expected += len(data[i])
			iov[i].Len = uint64(len(data[i]))
		}

		f, _ := conn.(*net.TCPConn).File()
		sz, err := vectorio.ReadvRaw(f.Fd(), iov)
		if err != nil {
			t.Errorf("read returned error %s", err.Error())
		}

		if sz != expected {
			t.Errorf("read was wrong length: %d != %d", sz, expected)
		}

		for i := range data {
			b := *(*[]byte)(unsafe.Pointer(&iov[i].Base))
			if bytes.Compare(data[i], b) != 0 {
				t.Errorf("read got wrong data: %s != %s", data[i], b)
			} else {
				t.Logf("read got correct data: %s == %s", data[i], b)
			}
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

	f, err := conn.File()
	if err != nil {
		t.Errorf("could not get file handle for TCP client: %s", err)
	}
	defer f.Close()
	nw, err := vectorio.Writev(f, data)
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
