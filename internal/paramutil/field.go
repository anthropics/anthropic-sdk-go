package paramutil

import "github.com/anthropics/anthropic-sdk-go/packages/param"

func AddrIfPresent[T comparable](v param.Opt[T]) *T {
	if v.IsPresent() {
		return &v.Value
	}
	return nil
}
