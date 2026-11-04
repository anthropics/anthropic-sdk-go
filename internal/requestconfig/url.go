package requestconfig

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Matches "." or ".." where each dot may be literal or percent-encoded
// (%2e / %2E). [url.PathEscape] leaves dots unchanged (they are unreserved
// characters), so a path parameter consisting only of dots would otherwise
// survive escaping and be resolved as a relative path segment, changing
// which route the request is sent to.
var dotSegmentRegexp = regexp.MustCompile(`^(?:\.|%2[eE]){1,2}$`)

// ValidateRequestPath returns an error if u, a URL or a path, fails [url.Parse]
// or would be sent with a path segment "." or ".." (a dot may also be %2e or
// %2E), since such a segment could change the route. When u is not fully
// escaped, Go re-escapes the decoded path: "%2e%2e%2f x" becomes "../%20x".
func ValidateRequestPath(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}
	segments := strings.Split(parsed.EscapedPath(), "/")
	if parsed.Scheme == "" {
		// requestURLReference trims the leading slashes of a relative
		// "//host/path", so its host becomes a path segment too.
		segments = append(segments, parsed.Host)
	}
	for _, segment := range segments {
		if dotSegmentRegexp.MatchString(segment) {
			return fmt.Errorf("constructed path %q contains dot segment %q which is not allowed", u, segment)
		}
	}
	return nil
}

// requestURLReference returns u in the form to resolve against the base URL.
// An absolute URL is unchanged. A relative one loses its leading slashes, so it
// resolves under the base URL's path, and gets a "./" prefix when its first
// segment has a colon (as in "/https://x/"), which Go would otherwise read as a
// scheme or reject (RFC 3986, section 4.2).
func requestURLReference(u *url.URL) string {
	if u.IsAbs() {
		return u.String()
	}
	ref := strings.TrimLeft(u.String(), "/")
	firstSegment := ref
	if i := strings.IndexAny(ref, "/?#"); i >= 0 {
		firstSegment = ref[:i]
	}
	if strings.Contains(firstSegment, ":") {
		return "./" + ref
	}
	return ref
}

// resolveRequestURL resolves req.URL against base, or sets it to nil on error.
// An absolute URL is kept as it is. A relative "//host/path" resolves under
// base, so req.Host is cleared when it names that host. A different req.Host
// is kept.
func resolveRequestURL(base *url.URL, req *http.Request) (err error) {
	if !req.URL.IsAbs() && req.Host == req.URL.Host {
		req.Host = ""
	}
	req.URL, err = base.Parse(requestURLReference(req.URL))
	return err
}
