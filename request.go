package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Header   map[string]interface{}
	Method   string
	Path     string
	Protocol string
	body     io.Reader
	Context  context.Context
}

func readRequest(reqReader *bufio.Reader, request *Request) error {

	// Read Header
	firstLine := true
	hasHeader := true
	for hasHeader {
		line, err := reqReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection (request canceled).")
				return nil
			}
			fmt.Println("Error: ", err)
			return err
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

	// Read Body
	if contentLength, ok := request.Header["Content-Length"]; ok {
		l, _ := strconv.Atoi(contentLength.(string))
		// Creating another buffer to read the body from the connection.
		buf := new(bytes.Buffer)

		// Copy from connection by reading and writing to the buffer.
		_, err := io.Copy(buf, io.LimitReader(reqReader, int64(l)))

		if err != nil {
			fmt.Println("Error reading 	body: ", err)
			return err // Client disconnected during body read.
		}
		// Assign the body to the request as reader body source.
		request.body = buf

	}

	return nil
}
