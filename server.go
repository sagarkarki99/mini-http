package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type HandleFunc func(w ResponseWriter, r *Request)

type Router struct {
	routes map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]HandleFunc)}
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
			fmt.Println("Listener error: ", err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			request := Request{}
			request.Header = make(map[string]interface{})
			hasHeader := true
			firstLine := true
			reqReader := bufio.NewReader(conn)
			for hasHeader {

				line, err := reqReader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						fmt.Println("Client closed the connection (request canceled).")
						return
					}
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

			responseW := &responseWriter{
				response: Response{},
			}

			if fn, ok := r.routes[request.Path]; ok {
				//create a context with deadline and cancel.
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				contentLength := request.Header["Content-Length"]
				cl, _ := strconv.Atoi(contentLength.(string))

				// Creating another buffer to read the body from the connection.
				buf := new(bytes.Buffer)

				// Copy from connection by reading and writing to the buffer.
				_, err := io.Copy(buf, io.LimitReader(reqReader, int64(cl)))

				if err != nil {
					fmt.Println("Error reading 	body: ", err)
					return // Client disconnected during body read.
				}
				// Assign the body to the request as reader body source.
				request.body = buf

				request.Context = ctx
				go monitorConnection(reqReader)

				// pass it to goroutine which checks if connection is active or not.
				// if connection is terminated, cancels the context.
				fn(responseW, &request)
				select {
				case <-ctx.Done():
					fmt.Println("Context canceled, skipping response write")
					return
				default:
					_, err := conn.Write([]byte("Hello world"))
					if err != nil {
						if err == net.ErrClosed {
							fmt.Println("Connection closed")
						} else {
							fmt.Println("Error writing:", err)
						}
					}
				}
			}
		}(con)
	}

}

func monitorConnection(reqReader *bufio.Reader) {
	// No need to check for context because , readByte is a blocking call from connection buffer.
	// if any operation happens (error, success), in connection, it will return.
	// This will return single byte from the buffer
	_, err := reqReader.ReadByte()
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "closed") {
			fmt.Println("Connection closed by client during processing.")
		} else {
			fmt.Println("Unexpected error while monitoring connection: ", err)
		}
		return
	}

}
