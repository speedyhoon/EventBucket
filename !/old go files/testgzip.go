package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"fmt"
)

func main () {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	input := []byte("hello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, worldhello, world")
	w.Write(input)
	w.Close() // You must close this first to flush the bytes to the buffer.
	fmt.Println(len( fmt.Sprintf("%v", w)), "\n\n", fmt.Sprintf("%v", w), "\n\n", fmt.Sprintf("%#v", w), "\n\n", fmt.Sprintf("%+v", w))
	fmt.Println("\n\n")
	fmt.Println("\n\n")
	fmt.Println(len( b.Bytes()), "\n\n", b.Bytes(), "\n\n", fmt.Sprintf("%v", b.Bytes()), "\n\n", fmt.Sprintf("%#v", b.Bytes()), "\n\n", bytes.IndexByte(b.Bytes(), 0), "\n\nLength: ", b.Len(), "\n\nLength: ", fmt.Sprintf("%+v", bytes.TrimSpace(b.Bytes())))
	err := ioutil.WriteFile("hello_world.txt.gz", b.Bytes(), 0666)
	if err != nil{
		fmt.Println(err)
	}

	if b.Len() < len(input) {
		fmt.Println("Cmpressed by: ", len(input) - b.Len())
	}

}
