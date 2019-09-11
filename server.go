package main

//go:generate go-assets-builder -s="/templates/" -o template.go -v Templates templates/
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	// "path/filepath"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

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

type Route struct {
	Paths []string
	Handler
}

func NewRoute(paths interface{}, handler Handler) *Route {
	switch v := paths.(type) {
	case string:
		return &Route{[]string{v}, handler}
	case []string:
		return &Route{v, handler}
	default:
		return &Route{nil, handler}
	}
}

type Router struct {
	Routes   []*Route
	NotFound *Route
	Error    *Route
}

func NewRouter() *Router {
	return &Router{
		NotFound: NewRoute(nil, func(req *Request) {
			res := Response{
				StatusCode: 404,
				StatusText: "Not Found",
				Headers:    map[string]string{},
				Body:       "Default Page Not Found",
				Conn:       req.Conn,
			}
			res.response()
		}),
		Error: NewRoute(nil, func(req *Request) {
			res := Response{
				StatusCode: 500,
				StatusText: "Internal Server Error",
				Headers:    map[string]string{},
				Body:       "Default Internal Server Error",
				Conn:       req.Conn,
			}
			res.response()
		}),
	}
}

func (router *Router) Match(path string) (*Route, error) {
	for _, route := range router.Routes {
		for _, _path := range route.Paths {
			if _path == path {
				return route, nil
			}
		}
	}
	return router.NotFound, nil
}

func (router *Router) Add(route *Route) {
	router.Routes = append(router.Routes, route)
}

func (router *Router) AddNotFound(route *Route) {
	router.NotFound = route
}

func (router *Router) AddError(route *Route) {
	router.Error = route
}

var (
	router *Router
)

func main() {

	router = &Router{}
	router.AddNotFound(NewRoute(nil, NotFoundHandler))

	router.Add(NewRoute("/api", ApiHandler))
	router.Add(NewRoute([]string{"", "/", "/home"}, HomeHandler))

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		req := translateHttp(&conn)

		fmt.Printf("[%s]: %s%s\n", time.Now().Format("2006/1/2 15:04:05"), req.Host, req.Path)
		if r, err := router.Match(req.Path); err == nil {
			go r.Handler(&req)
		}

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

	(*req.Conn).SetReadDeadline(time.Now().Add(100 * time.Millisecond))
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

func GetTemplate(template_name string) (string, error) {
	f, err := Templates.Open(template_name)
	if err != nil {
		return "", err
	} else {
		b, _ := ioutil.ReadAll(f)
		return string(b), nil
	}
}

func ApiHandler(req *Request) {
	req_info := map[string]interface{}{
		"method":      req.Method,
		"host":        req.Host,
		"source_host": req.Host,
		"path":        req.Path,
		"query":       req.Query,
		"data":        req.Data,
	}
	body, _ := json.MarshalIndent(req_info, "", "\t")
	res := Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		Conn:       req.Conn,
	}
	res.response()
}

func HomeHandler(req *Request) {
	html, _ := GetTemplate("home.html")
	res := Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       html,
		Conn:       req.Conn,
	}
	res.response()
}

func NotFoundHandler(req *Request) {
	html, _ := GetTemplate("404.html")
	res := Response{
		StatusCode: 404,
		StatusText: "Not Found",
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       html,
		Conn:       req.Conn,
	}
	res.response()
}

func ErrorHandler(req *Request) {
	html, _ := GetTemplate("500.html")
	res := Response{
		StatusCode: 500,
		StatusText: "Internal Server Error",
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       html,
		Conn:       req.Conn,
	}
	res.response()
}
