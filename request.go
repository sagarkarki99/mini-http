package main

import (
	"context"
	"io"
)

type Request struct {
	Header   map[string]interface{}
	Method   string
	Path     string
	Protocol string
	body     io.Reader
	Context  context.Context
}
