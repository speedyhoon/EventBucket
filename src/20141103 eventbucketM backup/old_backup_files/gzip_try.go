package main

import (
//	"net/http"
//	"mgo"
//	"strings"
//	"compress/gzip"
//	"io"
	"fmt"
	"bytes"
	"compress/gzip"
)

func main() {
	//BestCompression    = 9

	var b bytes.Buffer
//	 NewWriterLevel(w io.Writer, level int) (*Writer, error)
	gz, err := gzip.NewWriterLevel(&b, 9)
	defer gz.Close()
	_, err = gz.Write([]byte("YourDataHere"))
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}
