package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type HandleFunc func(w ResponseWriter, r *Request)

type Router struct {
	routes map[string]HandleFunc
}

func (r *Router) HandlerFunc(route string, fn HandleFunc) {
	if _, ok := r.routes[route]; ok {
		return
	}
	r.routes[route] = fn
}

func NewServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) StartListening(addr string, r *Router) error {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}
	for {
		con, err := lis.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go func(conn net.Conn) {
			request := Request{}
			request.Header = make(map[string]interface{})
			hasHeader := true
			firstLine := true
			reqReader := bufio.NewReader(conn)
			for hasHeader {
				line, err := reqReader.ReadString('\n')

				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				if firstLine {
					sep := strings.Split(line, " ")
					request.Method = strings.TrimSpace(sep[0])
					request.Path = strings.TrimSpace(sep[1])
					request.Protocol = strings.TrimSpace(sep[0])
					firstLine = false
				} else {
					if line == "\r\n" {
						break
					}
					sep := strings.Split(line, ":")
					key := strings.TrimSpace(sep[0])
					value := strings.TrimSpace(sep[1])
					request.Header[key] = value
				}
			}

			fmt.Println("Request: ", request)

			contentLength := request.Header["Content-Length"]
			cl, _ := strconv.Atoi(contentLength.(string))

			request.body = bufio.NewReaderSize(reqReader, cl)

			responseW := responseWriter{
				response: Response{},
			}

			if fn, ok := r.routes[request.Path]; ok {
				fn(responseW, &request)
			}

			conn.Write([]byte("Hello world"))
			conn.Close()
		}(con)
	}

}
