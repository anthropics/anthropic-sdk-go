package config

import "encoding/json"

// OAuthErrorBodyMaxLen caps the error body length embedded in OAuth-endpoint
// failure error messages. Token-endpoint responses are untrusted and are
// routinely captured in logs; the cap covers identity-provider proxies that
// concatenate upstream messages into error_description.
const OAuthErrorBodyMaxLen = 2000

// oauthErrorAllowedKeys is the RFC 6749 §5.2 allowlist surfaced in
// user-facing errors. Any other key is dropped — providers have been seen
// to echo back the assertion JWT or a partially-minted token, neither of
// which may appear in logs. error_description is attacker-controllable
// free-form text; callers must treat redacted output as untrusted.
var oauthErrorAllowedKeys = []string{
	"error",
	"error_description",
	"error_uri",
}

// RedactOAuthErrorBody returns a safe-to-log form of an OAuth token endpoint
// failure body. If the body parses as a JSON object, only the RFC 6749 §5.2
// allowed keys are kept. Non-JSON bodies are replaced with a redaction
// placeholder. Either way the result is truncated to [OAuthErrorBodyMaxLen].
func RedactOAuthErrorBody(body string) string {
	var obj map[string]any
	// err only — Unmarshal("null") sets obj=nil with no error; nil filters to {}.
	if err := json.Unmarshal([]byte(body), &obj); err == nil {
		filtered := make(map[string]any, len(oauthErrorAllowedKeys))
		for _, k := range oauthErrorAllowedKeys {
			if v, present := obj[k]; present {
				filtered[k] = v
			}
		}
		if b, err := json.Marshal(filtered); err == nil {
			body = string(b)
		}
	} else {
		body = "[redacted; not a JSON error response]"
	}
	if len(body) > OAuthErrorBodyMaxLen {
		body = body[:OAuthErrorBodyMaxLen] + "...[truncated]"
	}
	return body
}
