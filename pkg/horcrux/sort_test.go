package horcrux

import (
	"reflect"
	"sort"
	"testing"
)

func TestSortHumanReadable(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			"Sort files with multiple digits",
			[]string{"/a/b/c/test-1.text", "/a/b/c/test-10.text", "/a/b/c/test-3.text", "/a/b/c/test-100.text", "/a/b/c/test-11.text"},
			[]string{"/a/b/c/test-1.text", "/a/b/c/test-3.text", "/a/b/c/test-10.text", "/a/b/c/test-11.text", "/a/b/c/test-100.text"},
		},
		{
			"Sort files with multiple digits and multiple numbers",
			[]string{"1-1-30-1", "1-1-2-1", "1-1-1-4", "1-1-1-1"},
			[]string{"1-1-1-1", "1-1-1-4", "1-1-2-1", "1-1-30-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forceTerminal = true

			sort.Sort(NaturalSort(tt.input))

			if !reflect.DeepEqual(tt.input, tt.expect) {
				t.Errorf("SortHumanReadable() = %v, expected %v", tt.input, tt.expect)
			}
		})
	}
}
