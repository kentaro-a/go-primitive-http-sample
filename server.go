package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var (
	gcounter int = 0
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
	Conn       *net.Conn
}

type Request struct {
	Method     string
	Host       string
	SourceHost string
	Path       string
	Version    string
	Headers    map[string]string
	Data       map[string]string
	Query      map[string]string
	Conn       *net.Conn
}

func (res Response) response() {
	responses := []string{fmt.Sprintf("HTTP/1.1 %d %s", res.StatusCode, res.StatusText)}
	for k, v := range res.Headers {
		responses = append(responses, fmt.Sprintf("%s: %s", k, v))
	}
	responses = append(responses, "")
	responses = append(responses, res.Body)
	s := strings.Join(responses, "\r\n")
	(*res.Conn).Write([]byte(s))
	(*res.Conn).Close()
}

type Handler func(req *Request)
type Route string
type Router map[Route]Handler

var router Router

func main() {

	router = make(map[Route]Handler)

	router[Route("/")] = handler

	r1 := Route("/r1")
	router[r1] = handler1

	r2 := Route("/r2")
	router[r2] = handler2

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
		req := translateHttp(&conn)

		fmt.Printf("[%s]: %s%s\n", time.Now().Format("2006/1/2 15:04:05"), req.Host, req.Path)
		go router[Route(req.Path)](&req)
	}
}

func translateHttp(conn *net.Conn) Request {
	req := Request{
		Method:     "",
		Host:       (*conn).LocalAddr().String(),
		SourceHost: (*conn).RemoteAddr().String(),
		Path:       "",
		Version:    "",
		Headers:    map[string]string{},
		Data:       map[string]string{},
		Query:      map[string]string{},
		Conn:       conn,
	}

	(*req.Conn).SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	reader := bufio.NewReader(*req.Conn)
	line_num := 0
	is_body := false
	for {
		bb, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := string(bb)

		if line_num == 0 {
			parseRequestLine(line, &req)

		} else {
			if line == "" {
				is_body = true
			} else {
				if is_body {
					parseData(line, &req)
				} else {
					parseHeader(line, &req)
				}
			}
		}
		line_num++
	}
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

func parseRequestLine(str string, req *Request) {
	_elem := strings.Split(str, " ")
	req.Method = _elem[0]
	req.Path = strings.Split(_elem[1], "?")[0]
	req.Version = _elem[2]
	parseQuery(_elem[1], req)
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

func handler(req *Request) {
	req_info := map[string]interface{}{
		"method":      req.Method,
		"host":        req.Host,
		"source_host": req.Host,
		"path":        req.Path,
		"query":       req.Query,
		"data":        req.Data,
	}
	body, _ := json.Marshal(req_info)
	res := Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{"header1": "value1", "header2": "value2"},
		Body:       string(body),
		Conn:       req.Conn,
	}
	res.response()
}

func handler1(req *Request) {
	res := Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{},
		Body:       "handler1",
		Conn:       req.Conn,
	}
	res.response()
}

func handler2(req *Request) {
	res := Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{},
		Body:       "handler2",
		Conn:       req.Conn,
	}
	res.response()
}
