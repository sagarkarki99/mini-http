package main

import (
	"fmt"
	"net"
)

type ResponseWriter interface {
	WriteStatusCode(code int)
	Header(key string, value string)
	Write(data []byte)
}

type responseWriter struct {
	response Response
}

type Response struct {
	code   int
	header map[string]string
	data   []byte
}

func (w *responseWriter) WriteStatusCode(code int) {
	w.response.code = code
}

func (w *responseWriter) Header(key string, value string) {
	if w.response.header == nil {
		w.response.header = make(map[string]string)
	}
	w.response.header[key] = value
}

func (w *responseWriter) Write(data []byte) {
	w.response.data = data
}

// Add this new method to format the complete HTTP response
func (w *responseWriter) WriteToConnection(conn net.Conn) error {
	// Build the status line
	statusText := getStatusText(w.response.code)
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", w.response.code, statusText)

	// Write status line
	_, err := conn.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	// Write headers
	for key, value := range w.response.header {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := conn.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}

	// Write empty line to separate headers from body
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	// Write the response body
	if len(w.response.data) > 0 {
		_, err = conn.Write(w.response.data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to get status text for common HTTP status codes
func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}
