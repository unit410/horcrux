package horcrux

import "testing"

func TestPubkeyIsOmitted(t *testing.T) {
	r := Record{}
	if r.PubkeyIsOmitted() {
		t.Error("PubkeyIsOmitted(): Expected pubkey not to be omitted")
	}
	r.OmitPubkey()
	if !r.PubkeyIsOmitted() {
		t.Error("PubkeyIsOmitted(): Expected pubkey not to be omitted")
	}
}
