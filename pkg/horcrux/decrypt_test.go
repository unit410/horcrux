package horcrux

import "testing"

func Test_getEmailFromSmartcard(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"present", []byte("<foo@example.com>"), "foo@example.com"},
		{"present", []byte(`
Reader ...........: Yubico YubiKey FIDO CCID
Application ID ...: D1234567890123456789012345678901
Name of cardholder: foo@example.com
Language prefs ...: [not set]`), "foo@example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := getEmailFromSmartcard(tt.input); actual != tt.expected {
				t.Errorf("getEmailFromSmartcard() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
