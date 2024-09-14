package wallet

import (
	. "gochain/cryptor"
	"log"
	"os"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()

	t.Log(wallet.address)
}

func TestSignTransaction(t *testing.T) {
	wallet := NewWallet()

	signature, err := wallet.Sign("hello")
	if err != nil {
		t.Error(err)
	}
	t.Log(signature)

	if !wallet.Verify("hello", signature) {
		t.Error("Verification failed")
	}
	t.Log("Verification passed")
}

func TestSaveWallet(t *testing.T) {
	wallet := NewWallet()

	wallet.SaveToFile("wallet.dat")

	w, err := LoadWalletFromFile("wallet.dat")
	if err != nil {
		t.Error(err)
	}
	if w.address != wallet.address {
		t.Error("Address mismatch")
	}
	t.Log("Address match")
}

func TestSignAsString(t *testing.T) {
	wallet := NewWallet()

	pubString := wallet.GetPublicKeyAsString()
	t.Log(pubString)
	os.WriteFile("key.pem", []byte(pubString), 0644)

	signature, err := wallet.Sign("hello")
	if err != nil {
		t.Error(err)
	}
	log.Println(signature)

	valide, e := VerifySignatureWithPublicKeyString(pubString, "hello", signature)
	if e != nil {
		t.Error(e)
	}
	if !valide {
		t.Error("Verification failed")
	}
	t.Log("Verification passed")
}
