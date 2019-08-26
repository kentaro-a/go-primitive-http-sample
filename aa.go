package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	// "io"
	"net"
	"os"
	"strconv"
	"strings"
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
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Data    map[string]string
	Query   map[string]string
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
		// req := requestBinder(conn)
		// go handler(req)
		go requestBinder(conn)
	}
}

func requestBinder(conn net.Conn) Request {
	fmt.Print("Connected -> ")

	req := Request{
		Method:  "",
		Path:    "",
		Version: "",
		Headers: map[string]string{},
		Data:    map[string]string{},
		Query:   map[string]string{},
	}

	reader := bufio.NewReader(conn)
	line_num := 0

	for {
		bline, err := reader.ReadBytes('\n')
		fmt.Println(bline, err)
		if err == nil {
			line := string(bline[:len(bline)-2])
			fmt.Println(line_num, ": ", line)
			if line_num == 0 {
				_elem := strings.Split(line, " ")
				req.Method = _elem[0]
				req.Path = _elem[1]
				req.Version = _elem[2]
				parseQuery(req.Path, &req)

			} else {
				if line == "" {
					fmt.Println("[line]", line, ": empty")
				} else {
					fmt.Println("[line]", line, ": not empty")
				}
				parseHeader(line, &req)
			}
			line_num++

		} else {
			break
		}
	}
	fmt.Println("AAAAAAAAAAAA")

	if content_length, exist := req.Headers["Content-Length"]; exist {
		cl, err := strconv.Atoi(content_length)
		if err == nil {
			fmt.Println(cl)

		}
	}

	fmt.Println("Scanner end")
	return req
}

func parseData(str string, req *Request) {
	data := strings.Split(str, "&")
	for _, d := range data {
		s := strings.Split(d, "=")
		if len(s) == 2 && s[0] != "" {
			req.Data[s[0]] = s[1]
		}
	}
}

func parseQuery(str string, req *Request) {
	tmp := strings.Split(str, "?")
	if len(tmp) >= 2 {
		data := strings.Split(tmp[1], "&")
		for _, d := range data {
			s := strings.Split(d, "=")
			if len(s) == 2 && s[0] != "" {
				req.Query[s[0]] = s[1]
			}
		}
	}
}

func parseHeader(str string, req *Request) {
	headers := strings.Split(str, "&")
	for _, h := range headers {
		s := strings.Split(h, ":")
		if len(s) >= 2 && s[0] != "" {
			req.Headers[strings.TrimSpace(s[0])] = strings.TrimSpace(strings.Join(s[1:], ""))
		}
	}
}

func handler(conn net.Conn) {

}
