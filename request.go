package main

import "io"

type Request struct {
	Header   map[string]interface{}
	Method   string
	Path     string
	Protocol string
	body     io.Reader
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]HandleFunc)}
}
