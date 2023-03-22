package gosen

import (
	"net/http"
	"net/url"

	"github.com/manvalls/gosen/util"
)

func (h *wrappedHandler) sendEarlyHints(w http.ResponseWriter) {
	if h.runCache == nil {
		return
	}

	h.runCacheMux.RLock()

	for url := range h.runCache {
		w.Header().Add("Link", "<"+url+">; rel=preload; as=fetch")
	}

	h.runCacheMux.RUnlock()

	w.WriteHeader(http.StatusEarlyHints)
}

func (h *wrappedHandler) cacheRuns(version string, runList []string) {
	if h.runCache == nil {
		return
	}

	urls := []string{}

	query := "format=json"
	if version != "" {
		query += "&version=" + url.QueryEscape(version)
	}

	for _, run := range runList {
		urls = append(urls, util.AddToQuery(run, query))
	}

	somethingToChange := false

	h.runCacheMux.RLock()

	for _, url := range urls {
		if !h.runCache[url] {
			somethingToChange = true
			break
		}
	}

	h.runCacheMux.RUnlock()

	if somethingToChange {
		h.runCacheMux.Lock()

		for _, url := range urls {
			h.runCache[url] = true
		}

		h.runCacheMux.Unlock()
	}
}
