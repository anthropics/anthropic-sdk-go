package anthropic

// To accept the API's `string | []T` form, register string promotion
// for beta content unions. The generator already covers the non-beta
// ones (message.go init).

import (
	"reflect"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"

	"github.com/tidwall/gjson"
)

func registerStringPromotion[SliceT ~[]E, E any](wrap func(string) E) {
	apijson.RegisterCustomDecoder[SliceT](func(node gjson.Result, value reflect.Value, defaultDecoder func(gjson.Result, reflect.Value) error) error {
		if node.Type == gjson.String {
			arrayValue := reflect.MakeSlice(value.Type(), 1, 1)
			arrayValue.Index(0).Set(reflect.ValueOf(wrap(node.String())))
			value.Set(arrayValue)
			return nil
		}
		return defaultDecoder(node, value)
	})
}

func init() {
	registerStringPromotion[[]TextBlockParam](func(s string) TextBlockParam {
		return TextBlockParam{Text: s}
	})
	registerStringPromotion[[]BetaTextBlockParam](func(s string) BetaTextBlockParam {
		return BetaTextBlockParam{Text: s}
	})
	registerStringPromotion[[]BetaContentBlockParamUnion](func(s string) BetaContentBlockParamUnion {
		return BetaContentBlockParamUnion{OfText: &BetaTextBlockParam{Text: s}}
	})
	registerStringPromotion[[]BetaToolResultBlockParamContentUnion](func(s string) BetaToolResultBlockParamContentUnion {
		return BetaToolResultBlockParamContentUnion{OfText: &BetaTextBlockParam{Text: s}}
	})
}
