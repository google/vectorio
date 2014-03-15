vectorio: Scatter/Gather I/O for Go
===

vectorio is a library for Go (aka "golang") to implement the writev(2), and in the future, readv(2) system calls.  This enables a list of byte slices to be efficiently written to a file without first copying them to an intermediate buffer.  Go's built-in [bufio](http://golang.org/pkg/bufio/) library is the standard approach to buffered I/O, but vectorio can be a replacement to achieve the same goal but without copying data.

Example
---

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
