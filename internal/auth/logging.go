package auth

import (
	"log"
	"sync"
)

// warnOnce emits a log line for an error category exactly once per process.
// Used for best-effort failures that should not surface as request errors
// but also must not be swallowed silently (e.g. a permissions problem on
// the credentials cache would otherwise cause a fresh exchange on every
// request with no visible cause).
var (
	warnOnceMu  sync.Mutex
	warnOnceSet = map[string]bool{}
)

func warnOnce(key, format string, args ...any) {
	warnOnceMu.Lock()
	already := warnOnceSet[key]
	warnOnceSet[key] = true
	warnOnceMu.Unlock()
	if already {
		return
	}
	log.Printf("anthropic-sdk-go/auth: "+format, args...)
}

// ResetWarnOnceForTest clears the warnOnce dedupe set. Exported for test
// helpers in sibling packages (e.g. internal/auth integration tests in
// package auth_test) that need to observe a warning a prior test may
// already have triggered.
func ResetWarnOnceForTest() {
	warnOnceMu.Lock()
	warnOnceSet = map[string]bool{}
	warnOnceMu.Unlock()
}

// ConfigShadowSource describes where the static credential that is
// shadowing a profile came from. Used by [WarnConfigShadowed] to tailor
// the remediation hint.
type ConfigShadowSource int

const (
	// ConfigShadowFromEnv means the static credential was picked up from
	// ANTHROPIC_API_KEY or ANTHROPIC_AUTH_TOKEN by the env autoloader.
	// Remediation: unset the env var.
	ConfigShadowFromEnv ConfigShadowSource = iota
	// ConfigShadowFromOption means the static credential was passed
	// explicitly via option.WithAPIKey / option.WithAuthToken.
	// Remediation: remove the option.
	ConfigShadowFromOption
)

// WarnConfigShadowed warns once when a static credential is present
// alongside a profile config. Per the documented precedence, the static
// credential wins, which silently disables the profile's auth. Source
// determines whether the remediation points at an env var or an explicit
// option call.
func WarnConfigShadowed(name string, source ConfigShadowSource) {
	var remediation string
	switch source {
	case ConfigShadowFromOption:
		remediation = "remove the explicit option"
	default:
		remediation = "unset " + name + " (and remove any explicit WithAPIKey/WithAuthToken calls if present)"
	}
	warnOnce(
		"config-shadowed-by-"+name,
		"%s is set and takes precedence over the profile configuration passed via WithConfig; to use the profile, %s",
		name, remediation,
	)
}
