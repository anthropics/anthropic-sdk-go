package requestconfig

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// A dot segment is "." or "..", alone or followed by ";" and anything. A dot
// may be "%2e" or "%2E", and a ";" may be "%3B" or "%3b", as [url.PathEscape]
// writes it. PathEscape keeps dots, so ".." as a path parameter would change the
// route. Some servers drop a segment's ";" and the rest before they route, so
// "..;x" would too, and ";x" would become an empty segment.
var dotSegmentRegexp = regexp.MustCompile(`^(?:\.|%2[eE]){1,2}(?:$|;|%3[bB])`)

// ValidateRequestPath returns an error if u, a URL or a path, fails [url.Parse]
// or would be sent with a path segment that is "." or "..", or has nothing, "."
// or ".." before its first ";". Dots and ";" may be percent-encoded. When u is
// not fully escaped, Go escapes the decoded path: "%2e%2f x" becomes "./%20x".
func ValidateRequestPath(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}
	segments := strings.Split(parsed.EscapedPath(), "/")
	if parsed.Scheme == "" {
		// requestURLReference trims the leading slashes of a relative
		// "//user@host/path", so "user@host" becomes a path segment too.
		authority := parsed.Host
		if parsed.User != nil {
			authority = parsed.User.String() + "@" + authority
		}
		segments = append(segments, authority)
	}
	for _, segment := range segments {
		if (strings.HasPrefix(segment, ".") || strings.HasPrefix(segment, "%")) && dotSegmentRegexp.MatchString(segment) {
			return fmt.Errorf("constructed path %q contains dot segment %q which is not allowed", u, segment)
		}
		if strings.HasPrefix(segment, ";") || strings.HasPrefix(segment, "%3B") || strings.HasPrefix(segment, "%3b") {
			return fmt.Errorf("constructed path %q contains segment %q which starts with \";\" and is not allowed", u, segment)
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
