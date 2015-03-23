package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
)

var name *string = flag.String("n", "SB", "Your name")

func main() {
	flag.Parse()
	// fmt.Println(*name)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer conn.Close()

	conn.Write(bytes.NewBufferString(*name).Bytes())

	go readHandler(conn)

	for {
		msg := bufio.NewReader(os.Stdin)
		data, _, _ := msg.ReadLine()
		conn.Write(data)
	}

}

func readHandler(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Errorf("error")
			continue
		}
		fmt.Println(string(buf[:n]))

	}
}
