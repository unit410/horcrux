package horcrux

import (
	"io"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	var baseSlice = []string{"foo", "bar", "baz"}

	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{"contains",
			args{s: baseSlice, str: "foo"},
			true,
		},
		{
			"does not contain",
			args{s: baseSlice, str: "not present"},
			false,
		},
		{
			"empty",
			args{s: []string{}, str: "foo"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := Contains(tt.args.s, tt.args.str); actual != tt.expected {
				t.Errorf("Contains() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestAskForConfirmation(t *testing.T) {
	type args struct {
		source io.Reader
		prompt string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{"yes",
			args{source: strings.NewReader("yes\n")},
			true,
		},
		{"y",
			args{source: strings.NewReader("y\n")},
			true,
		},
		{"YES",
			args{source: strings.NewReader("YES\n")},
			true,
		},
		{"Y",
			args{source: strings.NewReader("Y\n")},
			true,
		},
		{"No",
			args{source: strings.NewReader("no\n")},
			false,
		},
		{"empty",
			args{source: strings.NewReader("\n")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := AskForConfirmation(tt.args.source, tt.args.prompt); actual != tt.expected {
				t.Errorf("AskForConfirmation() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
