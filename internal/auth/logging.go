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
	if firstWarning(key) {
		log.Printf("anthropic-sdk-go/auth: "+format, args...)
	}
}

func firstWarning(key string) bool {
	warnOnceMu.Lock()
	defer warnOnceMu.Unlock()
	already := warnOnceSet[key]
	warnOnceSet[key] = true
	return !already
}

// WarnAPIKeyAndAuthTokenInEnv warns once that the environment holds both an API
// key and an auth token, of which a client can send only the key.
func WarnAPIKeyAndAuthTokenInEnv() {
	if firstWarning("api-key-and-auth-token-in-env") {
		log.Print("Anthropic: ANTHROPIC_API_KEY and ANTHROPIC_AUTH_TOKEN are both set. Using ANTHROPIC_API_KEY.")
	}
}

// ResetWarnOnceForTest lets a test observe a warning that an earlier test has
// already triggered.
func ResetWarnOnceForTest() {
	warnOnceMu.Lock()
	warnOnceSet = map[string]bool{}
	warnOnceMu.Unlock()
}
