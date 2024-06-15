package main

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
)

func SolRestoreAccount(privateKey string) (account types.Account, publicKey string, err error) {
	account, err = types.AccountFromBase58(privateKey)
	if err != nil {
		return
	}
	publicKey = account.PublicKey.ToBase58()
	return
}

func SolVerifySign(publicKey string, message string, sig string) bool {
	publicKeyBytes, err := base58.Decode(publicKey)
	if err != nil {
		// log.Printf("publicKey decode error: %v\n", err)
		return false
	}
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		// log.Printf("sig decode error: %v\n", err)
		return false
	}
	return ed25519.Verify(publicKeyBytes, []byte(message), sigBytes)
}
