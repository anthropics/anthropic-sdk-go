package json

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"testing/synctest"
)

// A type encoder that is being built must be usable by a concurrent caller
// that is in a different synctest bubble, or in none. Sharing a bubbled
// sync.WaitGroup with such a caller makes Wait panic.
func TestTypeEncoderConcurrentFirstUseAcrossSynctestBubble(t *testing.T) {
	for i := 0; i < 500; i++ {
		typ := reflect.StructOf([]reflect.StructField{{
			Name: "F" + strconv.Itoa(i),
			Type: reflect.TypeOf(0),
			Tag:  reflect.StructTag(`json:"f"`),
		}})
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			typeEncoder(typ)(newEncodeState(), reflect.New(typ).Elem(), encOpts{})
		}()
		synctest.Test(t, func(t *testing.T) {
			close(start)
			typeEncoder(typ)(newEncodeState(), reflect.New(typ).Elem(), encOpts{})
		})
		wg.Wait()
	}
}
