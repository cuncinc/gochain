/*
 can refer to ethereum's transaction struct:
 https://ethereum.org/zh/developers/docs/transactions/
*/

package tx

import (
	"gochain/cryptor"
	"strconv"
	"time"
)

type Tx struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Value     int    `json:"value"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	PubKey    string `json:"pubKey"`
	Data      string `json:"data"`
}

/*
函数式选项模式只适合于复杂的实例化，如果参数只有简单几个，建议还是用普通的构造函数来解决。

	tran = NewTxWithOptions(
	WithFrom("bob"),
	WithTo("alice"),
	WithValue(34),
	WithSignature("bob's sig"),
	)
*/
type TxOptions func(t *Tx)

func NewTx(from string, to string, value int, pubKey string) *Tx {
	return &Tx{
		From:      from,
		To:        to,
		Value:     value,
		PubKey:    pubKey,
		Timestamp: time.Now().Unix(),
	}
}

func WithFrom(from string) TxOptions {
	return func(t *Tx) {
		t.From = from
	}
}

func WithTo(to string) TxOptions {
	return func(t *Tx) {
		t.To = to
	}
}

func WithValue(value int) TxOptions {
	return func(t *Tx) {
		t.Value = value
	}
}

func WithTimestamp(timestamp int64) TxOptions {
	return func(t *Tx) {
		t.Timestamp = timestamp
	}
}

func NewTxWithOptions(options ...TxOptions) *Tx {
	t := &Tx{}
	for _, option := range options {
		option(t)
	}
	return t
}

func (t *Tx) Validate() bool {
	valid, _ := cryptor.VerifySignatureWithPublicKeyString(t.PubKey, t.GetTxData(), t.Signature)
	return valid
}

func (t *Tx) GetTxData() string {
	return t.From + t.To + strconv.Itoa(t.Value) + strconv.FormatInt(t.Timestamp, 10)
}
