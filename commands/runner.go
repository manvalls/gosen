package commands

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Runner struct {
	GetRunHandler func(url string) http.Handler
	BaseRequest   *http.Request
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
	req.Header.Set("cookie", r.BaseRequest.Header.Get("cookie"))
	req.Header.Set("user-agent", r.BaseRequest.Header.Get("user-agent"))
	req.Header.Set("referer", r.BaseRequest.Header.Get("referer"))

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
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		routine.UnmarshalJSON(body)
		return
	}

}
