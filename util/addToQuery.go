package util

import "strings"

func AddToQuery(url string, value string) string {
	i := strings.Index(url, "?")
	if i == -1 {
		return url + "?" + value
	}

	if i == len(url)-1 {
		return url + value
	}

	return url + "&" + value
}
