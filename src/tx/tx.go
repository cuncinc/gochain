package tx

import (
	"gochain/cryptor"
	"strconv"
	"time"
)

type Tx struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int    `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	PubKey    string `json:"pub_key"`
}

/*
函数式选项模式只适合于复杂的实例化，如果参数只有简单几个，建议还是用普通的构造函数来解决。

	tran = NewTxWithOptions(
	WithFrom("bob"),
	WithTo("alice"),
	WithAmount(34),
	WithSignature("bob's sig"),
	)
*/
type TxOptions func(t *Tx)

func NewTx(from string, to string, amount int, pubKey string) *Tx {
	return &Tx{
		From:      from,
		To:        to,
		Amount:    amount,
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

func WithAmount(amount int) TxOptions {
	return func(t *Tx) {
		t.Amount = amount
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

func (t *Tx) Verify() bool {
	valid, _ := cryptor.VerifySignatureWithPublicKeyString(t.PubKey, t.GetTxData(), t.Signature)
	return valid
}

func (t *Tx) GetTxData() string {
	return t.From + t.To + strconv.Itoa(t.Amount) + strconv.FormatInt(t.Timestamp, 10)
}
