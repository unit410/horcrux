package horcrux

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

type KeyLocation int

const (
	KeyLocationUnknown KeyLocation = iota // 0
	KeyLocationLocal                      // 1
	KeyLocationStub                       // 2
)

// SecretKeyLocation identifies if the key is stored locally,
// is a stub (smartcard) or other
func SecretKeyLocation(keyID string) KeyLocation {
	// Refresh card-status to pick up stubs if a card is inserted
	smartcardIsAttached()

	// Now look at secret key location
	cmd := exec.Command("gpg", "--with-colons", "--list-secret-keys")
	stdout, stderr := cmd.Output()
	if stderr != nil {
		log.Fatal("Unable to list secret keys")
	}
	return parseSecretKeyLocation(stdout, keyID)
}

func parseSecretKeyLocation(rawStdout []byte, keyID string) KeyLocation {
	stdout := string(rawStdout)

	sc := bufio.NewScanner(strings.NewReader(stdout))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "sec:") || strings.HasPrefix(line, "ssb:") {
			components := strings.Split(line, ":")
			skID := components[4]
			if skID == keyID {
				// Field 15 - S/N of a token
				// '#': a stub on a card we haven't seen
				// 'D2760001240100000006086221730000': A card ID
				// '+': Secret key is available
				skType := components[14]
				if skType == "+" {
					return KeyLocationLocal
				} else if skType == "#" {
					// Stub is available
					return KeyLocationStub
				} else if len(skType) > 0 {
					// Stub is available and allocated to a smartcrad
					return KeyLocationStub
				}
			}
		}
	}
	return KeyLocationUnknown
}
