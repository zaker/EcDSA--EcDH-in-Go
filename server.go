package main

import "net"
import "fmt"

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		fmt.Printf("Received: %v", string(data))
		_, err = c.Write(data)
		if err != nil {
			panic(fmt.Sprint("Write: ", err))
		}
	}
}

// func main() {
// 	l, err := net.Listen("tcp", "0.0.0.0:5555")
// 	if err != nil {
// 		println("listen error", err)
// 		return
// 	}
//
// 	for {
// 		fd, err := l.Accept()
// 		if err != nil {
// 			println("accept error", err)
// 			return
// 		}
//
// 		go echoServer(fd)
// 	}
// }
