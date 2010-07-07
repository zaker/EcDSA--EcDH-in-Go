package main

import (
	"net"
	"os"
	"time"
)

var c net.Conn
var e os.Error

func main() {
	c, e = net.Dial("tcp", "", "localhost:5555")
	if e != nil {
		panic(e.String())
	}
	for {
		_, err := c.Write([]byte("hi\n"))
		if err != nil {
			println(err.String())
		}
		time.Sleep(1e9)
	}
}
