package block

import (
	"encoding/json"
	"gochain/tx"
	"gochain/wallet"
	"testing"
)

func TestBlock(t *testing.T) {
	alice := wallet.NewWallet()
	bob := wallet.NewWallet()

	tx1 := tx.NewTx(alice.GetAddress(), bob.GetAddress(), 9865, alice.GetPublicKeyAsString())
	tx1.Signature, _ = alice.Sign(tx1.Data)

	tx2 := tx.NewTx(bob.GetAddress(), alice.GetAddress(), 1234, bob.GetPublicKeyAsString())
	tx2.Signature, _ = bob.Sign(tx2.Data)

	transList := []tx.Tx{*tx1, *tx2}
	b := NewPoWBlock(1, "0", transList)
	b.MineBlock(4)
	str, _ := json.Marshal(b)
	t.Log(string(str))
	t.Log(b.Validate())
}
