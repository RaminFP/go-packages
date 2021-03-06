package main

import (
	"io"
	"os"
)

func main() {
	// Copy from stdin to stdout using io.CopyBuffer
	//
	// CopyBuffer is identical to Copy except that it stages through
	// the provided buffer (if one is required) rather than allocating
	// a temporary one. If buf is nil, one is allocated; otherwise
	// if it has zero length, CopyBuffer panics.
	io.CopyBuffer(os.Stdout, os.Stdin, nil)
}
