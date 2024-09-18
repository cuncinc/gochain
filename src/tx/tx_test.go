package tx

import (
	"encoding/json"
	. "gochain/wallet"
	"testing"
)

func TestNewTx(t *testing.T) {
	alice := NewWallet()
	bob := NewWallet()
	t1 := NewTx(alice.GetAddress(), bob.GetAddress(), 34, alice.GetPublicKeyAsString())
	sig, _ := alice.Sign(t1.GetTxData())
	t1.Signature = sig
	if !t1.Validate() {
		t.Error("签名验证失败")
	} else {
		t.Log("签名验证成功")
	}

	str, _ := json.Marshal(t1)
	t.Log(string(str))
}
