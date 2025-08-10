package main

type ResponseWriter interface {
	WriteStatusCode(code int)
	Header(key string, value string)
	Write(data []byte)
}

type responseWriter struct {
	response Response
}

type Response struct {
	code   int               `json:"code"`
	header map[string]string `json:"header"`
	data   []byte            `json:"data"`
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
