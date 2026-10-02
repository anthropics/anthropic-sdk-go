package requestconfig

import (
	"fmt"
	"net/url"
)

func encodePathParam(value string) string {
	switch value {
	case ".":
		return "%2E"
	case "..":
		return "%2E%2E"
	}
	return url.PathEscape(value)
}

// FormatPath escapes path parameters and inserts them into a request path.
//
// Each param is PathEscaped before interpolation so a literal "/" stays inside
// its own segment (foo/bar → foo%2Fbar) rather than becoming an extra path
// element. Bare "." / ".." are encoded so url.URL resolution cannot drop them.
func FormatPath(format string, params ...string) string {
	args := make([]any, len(params))
	for i, param := range params {
		args[i] = encodePathParam(param)
	}
	return fmt.Sprintf(format, args...)
}
