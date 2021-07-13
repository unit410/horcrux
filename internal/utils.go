package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func Assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// askForConfirmation asks the user for confirmation. A user must type "y" or "yes"
// Any other input will be considered as a No.
// Return true if the user confirmed with a yes, false otherwise.
func AskForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("%s [y/N]: ", prompt)

	response, err := reader.ReadString('\n')
	Assert(err)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

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
			out := string(stdout)
			re := regexp.MustCompile(`[\<^](.*?)[\^>]`)
			email := re.FindStringSubmatch(out)[1]
			if Contains(identities, email) {
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
			retry := AskForConfirmation("Failed to decrypt share. Retry?")
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
