package main

import "net"

var c net.Conn
var e error

// func main() {
// 	c, e = net.Dial("tcp", "localhost:5555")
// 	if e != nil {
// 		panic(e)
// 	}
// 	i := 0
// 	for {
// 		_, err := c.Write([]byte("hi" + strconv.Itoa(i) + "\n"))
// 		if err != nil {
// 			println(err)
// 		}
//
// 		time.Sleep(1e9)
// 		buf := make([]byte, 512)
// 		nr, err := c.Read(buf)
//
// 		data := buf[0:nr]
// 		fmt.Printf("echoed: %v", string(data))
// 		i++
// 	}
// }
