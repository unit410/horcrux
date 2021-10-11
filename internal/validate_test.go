package internal

import (
	"testing"

	"gitlab.com/unit410/vault-shamir/shamir"
)

type args struct {
	original  []byte
	parts     [][]byte
	threshold int
}

var original = []byte{3, 7, 9}

func TestCheckAllCombinations(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		shim     func([][]byte) [][]byte
		expected bool
	}{
		{
			"base 2/2",
			getArgs(t, original, 2, 2),
			func(i [][]byte) [][]byte {
				return i
			},
			true,
		},
		{
			"edited index 0 2/2",
			getArgs(t, original, 2, 2),
			func(i [][]byte) [][]byte {
				i[0][0] = 1
				return i
			},
			false,
		},
		{
			"edited index 1 2/2",
			getArgs(t, original, 2, 2),
			func(i [][]byte) [][]byte {
				i[1][0] = 1
				return i
			},
			false,
		},
		{
			"base 2/3",
			getArgs(t, original, 3, 2),
			func(i [][]byte) [][]byte {
				return i
			},
			true,
		},
		{
			"2/3 split with 3 given",
			args{original: original, parts: nil, threshold: 3},
			func(i [][]byte) [][]byte {
				ret, _ := shamir.Split(original, 3, 2)
				return ret
			},
			true,
		},
		{
			"edited index 0 2/3",
			getArgs(t, original, 3, 2),
			func(i [][]byte) [][]byte {
				i[0][0] = 1
				return i
			},
			false,
		},
		{
			"edited index 1 2/3",
			getArgs(t, original, 3, 2),
			func(i [][]byte) [][]byte {
				i[1][0] = 1
				return i
			},
			false,
		},
		{
			"edited index 2 2/3",
			getArgs(t, original, 3, 2),
			func(i [][]byte) [][]byte {
				i[2][0] = 1
				return i
			},
			false,
		},
		{
			"below threshold",
			args{original: original, parts: nil, threshold: 2},
			func(i [][]byte) [][]byte {
				ret, _ := shamir.Split(original, 3, 3)
				return ret
			},
			false,
		},
	}
	for _, tt := range tests {
		tt.args.parts = tt.shim(tt.args.parts)
		t.Run(tt.name, func(t *testing.T) {
			if actual := CheckAllCombinations(tt.args.original, tt.args.parts, tt.args.threshold); actual != tt.expected {
				t.Errorf("CheckAllCombinations() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

// getArgs builds expected args and error if shamir.Split fails
func getArgs(t *testing.T, original []byte, parts int, threshold int) args {
	split, err := shamir.Split(original, parts, threshold)
	if err != nil {
		t.Errorf("could not shamir split %s into %d/%d parts: %s", original, threshold, parts, err)
	}

	return args{
		original:  original,
		parts:     split,
		threshold: threshold,
	}
}
