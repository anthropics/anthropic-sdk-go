package json

import (
	"testing"
)

// Inner implements MarshalJSON to trigger the optimized code path
type benchInner struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (b benchInner) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{b.Name, b.Value})
}

// Nested structure with multiple MarshalJSON calls
type benchNested struct {
	Inner benchInner `json:"inner"`
	Items []int      `json:"items"`
}

func (b benchNested) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Inner benchInner `json:"inner"`
		Items []int      `json:"items"`
	}{b.Inner, b.Items})
}

// Deeply nested to amplify the effect
type benchDeep struct {
	Level1 benchNested `json:"level1"`
	Level2 benchNested `json:"level2"`
	Data   string      `json:"data"`
}

func (b benchDeep) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Level1 benchNested `json:"level1"`
		Level2 benchNested `json:"level2"`
		Data   string      `json:"data"`
	}{b.Level1, b.Level2, b.Data})
}

func BenchmarkMarshalNestedMarshalJSON(b *testing.B) {
	data := benchDeep{
		Level1: benchNested{
			Inner: benchInner{Name: "test1", Value: 100},
			Items: []int{1, 2, 3, 4, 5},
		},
		Level2: benchNested{
			Inner: benchInner{Name: "test2", Value: 200},
			Items: []int{6, 7, 8, 9, 10},
		},
		Data: "some test data here",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Slice of nested structs - common real-world pattern
func BenchmarkMarshalSliceOfNestedMarshalJSON(b *testing.B) {
	data := make([]benchDeep, 50)
	for i := range data {
		data[i] = benchDeep{
			Level1: benchNested{
				Inner: benchInner{Name: "test1", Value: i},
				Items: []int{1, 2, 3, 4, 5},
			},
			Level2: benchNested{
				Inner: benchInner{Name: "test2", Value: i * 2},
				Items: []int{6, 7, 8, 9, 10},
			},
			Data: "some test data here that is a bit longer to simulate real payloads",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
