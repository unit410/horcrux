package gpg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func GetEntities(files []string) ([]*openpgp.Entity, error) {

	var entities []*openpgp.Entity
	// read armored file and store into entites array
	// the size of this array will be the number of shares
	for _, file := range files {
		// we import into gpg for use with encryption
		cmd := exec.Command("gpg", "--import", file)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		entityList, err := getEntityListFromFile(file)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entityList...)
	}
	return entities, nil
}

func ReadEntity(pubkey []byte) (*openpgp.Entity, error) {
	packetReader := packet.NewReader(bytes.NewReader(pubkey))
	return openpgp.ReadEntity(packetReader)
}

// Import the given public key into the system's gpg keychain
func ImportPubkey(pubkey []byte) error {
	cmd := exec.Command("gpg", "--import")
	cmd.Stdin = bytes.NewBuffer(pubkey)
	err := cmd.Run()
	return err
}

// GetEncryptionPacketKeyID by listing packets and parsing out the
// key that this message was encrypted to
func GetEncryptionPacketKeyID(message []byte) (string, error) {
	cmd := exec.Command("gpg", "--list-packets", "--list-only")
	cmd.Stdin = bytes.NewBuffer(message)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	sc := bufio.NewScanner(strings.NewReader(string(stdout)))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, ":pubkey enc packet") {
			components := strings.Split(line, " ")
			keyID := components[len(components)-1]
			if isValidKeyID(keyID) {
				return keyID, nil
			}
		}
	}
	return "", errors.New("key not found")
}

// GetEntityListFromFile returns EntityList after reading it from the armor keyring file
func getEntityListFromFile(keyringFile string) (openpgp.EntityList, error) {
	keyringReader, err := os.Open(keyringFile)
	if err != nil {
		return nil, err
	}
	defer keyringReader.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(keyringReader)
	if err != nil {
		return nil, err
	}
	return entityList, nil
}

// SerializeWithoutSigs serializes the public part of the given Entity to w, excluding signatures from other entities
func SerializeWithoutSigs(entity *openpgp.Entity, w io.Writer) error {
	err := entity.PrimaryKey.Serialize(w)
	if err != nil {
		return err
	}
	for _, ident := range entity.Identities {
		err = ident.UserId.Serialize(w)
		if err != nil {
			return err
		}
		err = ident.SelfSignature.Serialize(w)
		if err != nil {
			return err
		}
	}
	for _, subkey := range entity.Subkeys {
		err = subkey.PublicKey.Serialize(w)
		if err != nil {
			return err
		}
		err = subkey.Sig.Serialize(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// Look up a partial keyID in all fingerprints int he provided entity,
// returning the full fingerprint or an empty string if not found
func FingerprintFromKeyID(entity *openpgp.Entity, keyID string) string {
	if !isValidKeyID(keyID) {
		return ""
	}

	primary_fp := fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint)
	if strings.HasSuffix(primary_fp, keyID) {
		return primary_fp
	}
	for _, sk := range entity.Subkeys {
		subkey_fp := fmt.Sprintf("%X", sk.PublicKey.Fingerprint)
		if strings.HasSuffix(subkey_fp, keyID) {
			return subkey_fp
		}
	}
	return ""
}

// Decrypt the given payload using the system's gpg keychain
func Decrypt(payload []byte) ([]byte, bytes.Buffer, error) {
	cmd := exec.Command("gpg", "--decrypt")

	var stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	return stdout, stderr, err
}

var keyIDRegex = regexp.MustCompile(`[0-9a-fA-F]{16}`)

func isValidKeyID(keyID string) bool {
	if len(keyID) != 16 {
		return false
	}
	loc := keyIDRegex.FindStringIndex(keyID)
	if len(loc) == 0 || loc[0] != 0 {
		return false
	}
	return true
}
