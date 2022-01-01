package horcrux

import (
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
)

// getEntityListFromFile returns EntityList after reading it from the armor keyring file
func getEntityListFromFile(keyFile string) openpgp.EntityList {
	keyringReader, err := os.Open(keyFile)
	Assert(err)
	defer keyringReader.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(keyringReader)
	Assert(err)
	return entityList
}

// serializeWithoutSigs serializes the public part of the given Entity to w, excluding signatures from other entities
func serializeWithoutSigs(entity *openpgp.Entity, w io.Writer) error {
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
