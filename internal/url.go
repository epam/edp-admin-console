package internal

import "strings"

const (
	HttpsScheme = "https://"
	HttpScheme  = "http://"
)

func AddSchemeIfNeeded(url string) string {
	if !strings.HasPrefix(url, HttpScheme) && !strings.HasPrefix(url, HttpsScheme) {
		url = HttpsScheme + url
	}
	return url
}
