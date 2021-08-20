package internal

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

// DecryptPayload decrypts and returns the payload encrypted to pubkey
func DecryptPayload(payload []byte, pubkey []byte) (share []byte) {
	packetReader := packet.NewReader(bytes.NewReader(pubkey))
	entity, err := openpgp.ReadEntity(packetReader)
	Assert(err)

	fingerprint := entity.PrimaryKey.Fingerprint

	log.Printf("----------------------------------------------------------------\n")
	log.Printf("%X\n", fingerprint)
	identities := make([]string, 0, len(entity.Identities))
	for _, value := range entity.Identities {
		userId := value.UserId
		log.Printf("%s\n", userId.Name)
		log.Printf("%s\n", userId.Email)
		identities = append(identities, userId.Email)
	}

	{
		// we must import the pubkey before gpg --card-status can pair with the private key
		cmd := exec.Command("gpg", "--import")
		cmd.Stdin = bytes.NewBuffer(pubkey)
		err = cmd.Run()
		Assert(err)
	}

	log.Printf("Waiting for the above identity's smartcard to be inserted")

	for {
		fmt.Print(".")
		time.Sleep(300 * time.Millisecond)
		cmd := exec.Command("gpg", "--card-status")
		stdout, stderr := cmd.Output()
		if stderr == nil {
			if Contains(identities, getEmailFromSmartcard(stdout)) {
				log.Printf("\nSmartcard detected...\n")
				break
			}
		}
	}

	for {
		// ask gpg to decrypt the file
		log.Printf("Decrypting %x share...\n", fingerprint)
		cmd := exec.Command("gpg", "--decrypt")

		var stderr bytes.Buffer
		cmd.Stdin = bytes.NewReader(payload)
		cmd.Stderr = &stderr
		stdout, err := cmd.Output()
		if err != nil {
			log.Printf("%s\n", stderr.String())
			retry := AskForConfirmation(os.Stdin, "Failed to decrypt share. Retry?")
			if !retry {
				break
			}

			// try again
			continue
		}

		log.Printf("%x share decrypted with size %d\n", fingerprint, len(stdout))
		return stdout
	}

	return nil
}

// getEmailFromSmartcard decodes gpg --card-status output and finds the email address
// associated with the given smart card.
func getEmailFromSmartcard(input []byte) string {
	out := string(input)
	re := regexp.MustCompile(`[\<^](.*?)[\^>]`)
	return re.FindStringSubmatch(out)[1]
}
