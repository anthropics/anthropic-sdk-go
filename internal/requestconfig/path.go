package requestconfig

import (
	"fmt"
	"net/url"
	"strings"
)

// Generated resource methods append this exact query suffix; preserve it so
// PathEscape of ID segments cannot swallow the SDK's own beta flag.
const generatedBetaQuery = "?beta=true"

// escapeRequestPath PathEscapes each slash-separated segment of a relative
// request path so caller-controlled IDs cannot inject "?", "#", or extra
// query/fragment delimiters. Remaining "." / ".." segments (including
// percent-encoded forms) are rejected so url.URL.Parse / ResolveReference
// cannot walk out of the intended prefix.
func escapeRequestPath(u string) (string, error) {
	query := ""
	path := u
	if strings.HasSuffix(path, generatedBetaQuery) {
		path = strings.TrimSuffix(path, generatedBetaQuery)
		query = generatedBetaQuery
	}

	segments := strings.Split(path, "/")
	for i, seg := range segments {
		decoded := seg
		if d, err := url.PathUnescape(seg); err == nil {
			decoded = d
		}
		if decoded == "." || decoded == ".." {
			return "", fmt.Errorf("requestconfig: invalid path segment %q", seg)
		}
		segments[i] = url.PathEscape(seg)
	}
	return strings.Join(segments, "/") + query, nil
}
