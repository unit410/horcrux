package horcrux

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

// DecryptPayload decrypts and returns the payload encrypted to pubkey
func DecryptPayload(payload []byte, pubkey []byte) (share []byte) {
	packetReader := packet.NewReader(bytes.NewReader(pubkey))
	entity, err := openpgp.ReadEntity(packetReader)
	Assert(err)

	fingerprint := entity.PrimaryKey.Fingerprint

	logf("----------------------------------------------------------------\n")
	logf("%X\n", fingerprint)
	identities := make([]string, 0, len(entity.Identities))
	for _, value := range entity.Identities {
		userId := value.UserId
		logf("Name: %s\n", userId.Name)
		logf("Email: %s\n", userId.Email)
		identities = append(identities, userId.Email)
	}

	{
		// we must import the pubkey before gpg --card-status can pair with the private key
		cmd := exec.Command("gpg", "--import")
		cmd.Stdin = bytes.NewBuffer(pubkey)
		err = cmd.Run()
		Assert(err)
	}

	logf("Waiting for the above identity's smartcard to be inserted")

	for {
		fmt.Fprintf(os.Stderr, ".")
		time.Sleep(300 * time.Millisecond)
		cmd := exec.Command("gpg", "--card-status")
		stdout, stderr := cmd.Output()
		if stderr == nil {
			if Contains(identities, getEmailFromSmartcard(stdout)) {
				logf("\nSmartcard detected...\n")
				break
			}
		}
	}

	for {
		// ask gpg to decrypt the file
		logf("Decrypting %x share...\n", fingerprint)
		cmd := exec.Command("gpg", "--decrypt")

		var stderr bytes.Buffer
		cmd.Stdin = bytes.NewReader(payload)
		cmd.Stderr = &stderr
		stdout, err := cmd.Output()
		if err != nil {
			logf("%s\n", stderr.String())
			retry := AskForConfirmation(os.Stdin, "Failed to decrypt share. Retry?")
			if !retry {
				break
			}

			// try again
			continue
		}

		logf("%x share decrypted with size %d\n", fingerprint, len(stdout))
		return stdout
	}

	return nil
}

// getEmailFromSmartcard decodes gpg --card-status output and finds the email address
// associated with the given smart card.
func getEmailFromSmartcard(input []byte) string {
	out := string(input)

	// Canonical Form
	re := regexp.MustCompile(`[\<^](.*?)[\^>]`)
	matches := re.FindStringSubmatch(out)
	if len(matches) >= 2 {
		return matches[1]
	}
	// Fallback to Cardholder Name
	re = regexp.MustCompile(`Name of cardholder:\s+(.*)`)
	matches = re.FindStringSubmatch(out)
	if len(matches) >= 2 {
		logln("Warning: Could not find a canonical email, using full name as email")
		return matches[1]
	}
	logln("Warning: Could not find an email associated with smartcard")
	return ""
}
