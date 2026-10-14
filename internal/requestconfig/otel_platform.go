package requestconfig

// WithOpenTelemetryPlatform returns a RequestOption for the client options of a
// third-party platform: the instrumentation scope of the OpenTelemetry spans is
// that of the SDK with platform, the name of the client's package, after it.
func WithOpenTelemetryPlatform(platform string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.OpenTelemetry.Platform = platform
		return nil
	})
}
