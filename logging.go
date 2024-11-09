package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type incomingTmplValues struct {
	TraceID     string
	Method      string
	URL         string
	Protocol    string
	Accept      string
	Host        string
	ContentType string
}

var incomingTmpl = template.Must(template.New("").Parse(`
Incoming Request: {{ .TraceID }}
{{ .Method }} {{ .URL }} {{ .Protocol }}
Accept: {{ .Accept }}
Host: {{ .Host }}
Content-Type: {{ .ContentType }}
`))

func fromRequest(r *http.Request) *incomingTmplValues {
	return &incomingTmplValues{
		TraceID:     r.Header.Get("Traceparent"),
		Method:      r.Method,
		URL:         r.URL.String(),
		Protocol:    r.Proto,
		Accept:      r.Header.Get("Accept"),
		Host:        r.Host,
		ContentType: r.Header.Get("Content-Type"),
	}
}

func logRequest(r *http.Request) {
	err := incomingTmpl.Execute(os.Stdout, fromRequest(r))
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return
	}
}

var _ http.ResponseWriter = (*ResponseWriterWrapper)(nil)
var _ http.Flusher = (*ResponseWriterWrapper)(nil)
var _ http.Hijacker = (*ResponseWriterWrapper)(nil)

type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	f          *http.Flusher
	h          *http.Hijacker
	statusCode *int
}

func (rww ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return (*rww.h).Hijack()
}

func (rww ResponseWriterWrapper) Flush() {
	(*rww.f).Flush()
}

func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

func (rww ResponseWriterWrapper) Write(bytes []byte) (int, error) {
	return (*rww.w).Write(bytes)
}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	*rww.statusCode = statusCode
	(*rww.w).WriteHeader(statusCode)
}

// NewResponseWriterWrapper static function creates a wrapper for the http.ResponseWriter
func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var statusCode int = 200
	// Every request should implement flusher
	flusher, _ := w.(http.Flusher)
	// Every request should implement hijacker
	hijacker, _ := w.(http.Hijacker)
	return ResponseWriterWrapper{
		w:          &w,
		h:          &hijacker,
		f:          &flusher,
		statusCode: &statusCode,
	}
}

type outgoingTmplValues struct {
	TraceID     string
	Duration    string
	Protocol    string
	StatusCode  int
	ContentType string
}

var outgoingTmpl = template.Must(template.New("").Parse(`
Outgoing Response: {{ .TraceID }}
Duration: {{ .Duration }}
{{ .Protocol }} {{ .StatusCode }}
Content-Type: {{ .ContentType }}
`))

func fromResponse(w *ResponseWriterWrapper, r *http.Request, duration time.Duration) *outgoingTmplValues {
	return &outgoingTmplValues{
		TraceID:     (*w.w).Header().Get("Traceparent"),
		Duration:    duration.String(), // To seconds
		Protocol:    r.Proto,
		StatusCode:  *w.statusCode,
		ContentType: (*w.w).Header().Get("Content-Type"),
	}
}

func logResponse(w *ResponseWriterWrapper, r *http.Request, duration time.Duration) {
	err := outgoingTmpl.Execute(os.Stdout, fromResponse(w, r, duration))
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return
	}
}

func HTTPLogging(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique trace ID for the request
		traceID := uuid.New().String()

		// TODO: Actually add open-telemetry compatibility
		r.Header.Set("Traceparent", traceID)
		logRequest(r)
		rww := NewResponseWriterWrapper(w)
		start := time.Now()
		next(rww, r)
		duration := time.Since(start)
		logResponse(&rww, r, duration)
	})
}
