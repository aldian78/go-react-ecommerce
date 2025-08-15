package model

import (
	"bufio"
	"bytes"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

// ResponseBuffer .
type ResponseBuffer struct {
	Response gin.ResponseWriter // the actual ResponseWriter to flush to
	status   int                // the HTTP response code from WriteHeader
	Body     *bytes.Buffer      // the response content body
	Flushed  bool
}

// NewResponseBuffer .
func NewResponseBuffer(w gin.ResponseWriter) *ResponseBuffer {
	return &ResponseBuffer{
		Response: w, status: defaultStatus, Body: &bytes.Buffer{},
	}
}

// Header .
func (w *ResponseBuffer) Header() http.Header {
	return w.Response.Header() // use the actual response header
}

func (w *ResponseBuffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

// WriteString .
func (w *ResponseBuffer) WriteString(s string) (n int, err error) {
	//w.WriteHeaderNow()
	//n, err = io.WriteString(w.ResponseWriter, s)
	//w.size += n
	n, err = w.Write([]byte(s))
	return
}

// Written .
func (w *ResponseBuffer) Written() bool {
	return w.Body.Len() != noWritten
}

// WriteHeader .
func (w *ResponseBuffer) WriteHeader(status int) {
	w.status = status
}

// WriteHeaderNow .
func (w *ResponseBuffer) WriteHeaderNow() {
	//if !w.Written() {
	//	w.size = 0
	//	w.ResponseWriter.WriteHeader(w.status)
	//}
}

// Status .
func (w *ResponseBuffer) Status() int {
	return w.status
}

// Size .
func (w *ResponseBuffer) Size() int {
	return w.Body.Len()
}

// Hijack .
func (w *ResponseBuffer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	//if w.size < 0 {
	//	w.size = 0
	//}
	return w.Response.(http.Hijacker).Hijack()
}

// CloseNotify .
func (w *ResponseBuffer) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

// Flush .
func (w *ResponseBuffer) Flush() {
	w.realFlush()
}

func (w *ResponseBuffer) realFlush() {
	if w.Flushed {
		return
	}
	w.Response.WriteHeader(w.status)
	if w.Body.Len() > 0 {
		_, err := w.Response.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}
	w.Flushed = true
}

// Pusher .
func (w *ResponseBuffer) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.Response.(http.Pusher); ok {
		return pusher
	}
	return nil
}
