package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	// "time"
)

var m map[string]*net.Conn

func main() {

	m = make(map[string]*net.Conn, 10)
	index := 0
	var name string

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer listener.Close()

	for {
		//接收并创建用户；
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf(err.Error())
			return
		}
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)

		if name = string(buf[:n]); name == "SB" {
			name = name + strconv.Itoa(index)
			index++
		}

		// name := string(buf)

		fmt.Println(name, "正在连接")

		m[name] = &conn
		go func() {
			defer connCloseHandle(name, conn)
			connHandler(name, conn)
		}()
	}
}

func connCloseHandle(name string, conn net.Conn) {
	conn.Close()
	delete(m, name)
	for na, _ := range m {
		if na != name {
			//注意一下表达，(*buffer).bytes()
			// (*(m[na])).Write(bytes.NewBufferString(name + "离开").Bytes())
			// (*(m[na])).Write(bytes.NewBufferString("当前剩余" + strconv.Itoa(len(m)) + "人").Bytes())
			(*(m[na])).Write(bytes.NewBufferString(name + "离开").Bytes())
			(*(m[na])).Write(bytes.NewBufferString("当前剩余" + strconv.Itoa(len(m)) + "人").Bytes())
		}
	}
}

func connHandler(name string, conn net.Conn) {
	conn.Write(bytes.NewBufferString("当前在线" + strconv.Itoa(len(m)) + "人").Bytes())
	for na, _ := range m {
		if na != name {
			(*(m[na])).Write(bytes.NewBufferString(name + "进入房间").Bytes())
		}
	}
	for {
		//此处长度。。可能要修改
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		msg := string(buf[:n])
		conn.Write(bytes.NewBufferString("我说:" + msg).Bytes())
		for na, _ := range m {
			if na != name {
				(*(m[na])).Write(bytes.NewBufferString(name + "说：" + msg).Bytes())
			}
		}
	}
}
