package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type Size int

func (s Size) String() string {
	const unit = 1024
	if s < unit {
		return fmt.Sprintf("%d B", s)
	}
	div, exp := int64(unit), 0
	for n := s / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(s)/float64(div), "KMGTPE"[exp])
}

type SizeRecorder struct {
	http.ResponseWriter

	Size   Size
	Status int
}

func (sr *SizeRecorder) WriteHeader(code int) {
	sr.ResponseWriter.WriteHeader(code)
	sr.Status = code
}

func (sr *SizeRecorder) Write(b []byte) (n int, err error) {
	n, err = sr.ResponseWriter.Write(b)
	sr.Size = sr.Size + Size(n)

	return n, err
}

func (sr *SizeRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w, ok := sr.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return w.Hijack()
}

func (sr *SizeRecorder) Flush() {
	w, ok := sr.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	w.Flush()
}

type SizeRequest struct {
	Source io.ReadCloser
	Size   Size
}

func (sr *SizeRequest) Read(p []byte) (n int, err error) {
	n, err = sr.Source.Read(p)
	if err != nil {
		return n, err
	}
	sr.Size = sr.Size + Size(n)

	return n, nil
}

func (sr *SizeRequest) Close() error {
	return sr.Source.Close()
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		srec := &SizeRecorder{
			ResponseWriter: rw,
			Status:         200,
			Size:           0,
		}

		const format = "| %-3s | %3d | %-7s | In: %15s | Out: %15s | %10s | %45s | %s"
		var (
			start    = time.Now()
			uri      = r.RequestURI
			method   = r.Method
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		)

		sreq := &SizeRequest{Source: r.Body}
		r.Body = sreq
		log.Printf(
			format,
			"IN",
			0,
			r.Method,
			"?",
			"?",
			"?",
			ip,
			uri,
		)

		h.ServeHTTP(srec, r) // serve the original request

		duration := time.Since(start).Round(time.Millisecond)

		// round duration to 1 digit after seconds, if the duration is bigger than 1 second
		if duration > 1*time.Second {
			// A duration are Nanoseconds as int64
			// Milliseconds are 1000 Microseconds, which are 1000 Nanoseconds.
			// To round it up/down, we need to add 50 Milliseconds. Or 50,000,000 Nanoseconds.
			duration = ((duration + 5e7) / 1e8) * 1e8
		}

		// log request details
		log.Printf(
			format,
			"OUT",
			srec.Status,
			method,
			sreq.Size,
			srec.Size,
			duration,
			ip,
			uri,
		)
	}
	return http.HandlerFunc(logFn)
}
