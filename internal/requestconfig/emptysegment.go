package requestconfig

import "strings"

// IsEmptySegment reports whether segment, a path segment before it is escaped,
// is empty or has only ";" parameters. Some servers drop the parameters before
// routing, so they read "/a/;x/b" as "/a//b".
func IsEmptySegment(segment string) bool {
	name, _, _ := strings.Cut(segment, ";")
	return name == ""
}
