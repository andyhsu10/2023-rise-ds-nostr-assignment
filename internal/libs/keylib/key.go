package keylib

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"

	"distrise/internal/configs"
)

func Generate() {
	sk := nostr.GeneratePrivateKey()
	pk, _ := nostr.GetPublicKey(sk)
	nsec, _ := nip19.EncodePrivateKey(sk)
	npub, _ := nip19.EncodePublicKey(pk)

	fmt.Println("sk:", sk)
	fmt.Println("pk:", pk)
	fmt.Println(nsec)
	fmt.Println(npub)
}

func GetKeyPair() Key {
	config := configs.GetConfig()
	nsec, _ := nip19.EncodePrivateKey(config.PrivateKey)
	npub, _ := nip19.EncodePublicKey(config.PublicKey)

	return Key{
		Public:  npub,
		Private: nsec,
	}
}

func GetNPub(pk string) string {
	npub, _ := nip19.EncodePublicKey(pk)
	return npub
}

func GetNSec(sk string) string {
	nsec, _ := nip19.EncodePrivateKey(sk)
	return nsec
}

type Key struct {
	Public  string
	Private string
}
