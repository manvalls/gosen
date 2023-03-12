package commands

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type Runner struct {
	GetRunHandler func(url string) http.Handler
	BaseRequest   *http.Request
	Header        http.Header
}

type RunnerWriter struct {
	*Routine

	writer     io.Writer
	header     http.Header
	statusCode int
}

func (w *RunnerWriter) Header() http.Header {
	return w.header
}

func (w *RunnerWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *RunnerWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (r *Runner) Run(routine *Routine, url string) {
	req, err := http.NewRequestWithContext(r.BaseRequest.Context(), "GET", url, strings.NewReader(""))
	if err != nil {
		return
	}

	if req.Host == "" {
		req.Host = r.BaseRequest.Host
	}

	if req.Method == "" {
		req.Method = "GET"
	}

	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}

	if req.URL.Host == "" {
		req.URL.Host = r.BaseRequest.Host
	}

	req.Header.Set("gosen-accept", "json")
	req.Header.Set("accept-language", r.BaseRequest.Header.Get("accept-language"))
	req.Header.Set("user-agent", r.BaseRequest.Header.Get("user-agent"))
	req.Header.Set("referer", r.BaseRequest.URL.String())

	req.Header.Set("x-forwarded-for", r.BaseRequest.Header.Get("x-forwarded-for"))
	req.Header.Set("x-forwarded-proto", r.BaseRequest.Header.Get("x-forwarded-proto"))
	req.Header.Set("x-forwarded-host", r.BaseRequest.Header.Get("x-forwarded-host"))
	req.Header.Set("x-forwarded-port", r.BaseRequest.Header.Get("x-forwarded-port"))
	req.Header.Set("x-real-ip", r.BaseRequest.Header.Get("x-real-ip"))

	handler := r.GetRunHandler(url)
	if handler == nil {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}

		routine.UnmarshalJSON(body)
		return
	}

	buff := new(bytes.Buffer)
	rw := &RunnerWriter{routine, buff, make(http.Header), 200}
	handler.ServeHTTP(rw, req)

	if rw.statusCode/100 == 3 {
		r.Run(routine, rw.header.Get("Location"))
		return
	}

	body := buff.Bytes()
	if len(body) > 0 && body[0] == '[' && body[len(body)-1] == ']' {
		routine.UnmarshalJSON(body)
	}
}