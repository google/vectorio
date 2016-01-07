vectorio: Scatter/Gather I/O for Go
===

vectorio is a library for Go (aka "golang") to implement the writev(2), and in the future, readv(2) system calls.  This enables a list of byte slices to be efficiently written to a file without first copying them to an intermediate buffer.  Go's built-in [bufio](http://golang.org/pkg/bufio/) library is the standard approach to buffered I/O, but vectorio can be a replacement to achieve the same goal but without copying data.

This is not an official Google product.  Authored by Ted Behling as a personal project.

[![Build Status](https://travis-ci.org/google/vectorio.svg)](https://travis-ci.org/google/vectorio)

See [example usage](vectorio_example_test.go) for a demo.

