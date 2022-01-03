package horcrux

import (
	"reflect"
	"testing"
)

func TestParseSmartcardFingerprints(t *testing.T) {
	tests := []struct {
		name     string
		stdout   string
		expected []string
	}{
		{
			"Parse Fingerprints",
			`
Reader:Yubico YubiKey FIDO CCID:AID:D2760000000:openpgp-card:
version:0000:
vendor:0000:Yubico:
serial:00000000:
name:First:Last:
lang::
sex:m:
url::
login::
forcepin:0:::
keyattr:1:1:4096:
keyattr:2:1:4096:
keyattr:3:1:4096:
maxpinlen:127:127:127:
pinretry:3:0:3:
sigcount:0:::
uif:0:1:1:
cafpr::::
fpr:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA:BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB:CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC:
fprtime:0000000000:0000000000:0000000000:
grp:1111111111111111111111111111111111111111:2222222222222222222222222222222222222222:3333333333333333333333333333333333333333:
	`,
			[]string{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB", "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"},
		},
		{
			"Parse Missing Fingerprints",
			`
Reader:Yubico YubiKey FIDO CCID:AID:D2760000000:openpgp-card:
version:0000:
vendor:0000:Yubico:
serial:00000000:
name:First:Last:
lang::
sex:m:
url::
login::
forcepin:0:::
keyattr:1:1:4096:
keyattr:2:1:4096:
keyattr:3:1:4096:
maxpinlen:127:127:127:
pinretry:3:0:3:
sigcount:0:::
uif:0:1:1:
cafpr::::
fpr:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA::CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC:
fprtime:0000000000:0000000000:0000000000:
grp:1111111111111111111111111111111111111111:2222222222222222222222222222222222222222:3333333333333333333333333333333333333333:
	`,
			[]string{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseSmartcardFingerprints([]byte(tt.stdout))
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("parseSmartcardFingerprints() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
