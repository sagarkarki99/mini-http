package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
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
			request := Request{
				Header: make(map[string]interface{}),
			}
			reqReader := bufio.NewReader(conn)
			err := readRequest(reqReader, &request)
			if err != nil {
				fmt.Println("Error reading request: ", err)
				return // Client disconnected during request read.
			}

			if fn, ok := r.routes[request.Path]; ok {
				//create a context with deadline and cancel.
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				request.Context = ctx
				go monitorConnection(reqReader, cancel)

				// pass it to goroutine which checks if connection is active or not.
				// if connection is terminated, cancels the context.
				responseW := &responseWriter{
					response: Response{},
				}
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

func monitorConnection(reqReader *bufio.Reader, cancel context.CancelFunc) {
	defer cancel()
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
