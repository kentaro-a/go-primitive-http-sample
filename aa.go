package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	// "reflect"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Text       string `json:"text"`
	Conn       net.Conn
}

func (res Response) response() {
	b, _ := json.Marshal(res)
	res.Conn.Write(b)
	res.Conn.Write([]byte("\n"))
	res.Conn.Close()
}

type Request struct {
	Data   int    `json:"status_code"`
	Method string `json:"text"`
}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var gcounter int = 0

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := l.Accept()
		gcounter++
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		res := Response{
			StatusCode: 200,
			Text:       "Ok",
			Conn:       conn,
		}
		res.response()
	}
}
