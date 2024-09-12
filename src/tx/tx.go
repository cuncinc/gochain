package tx

import "strings"

type Tx struct {
	From      string
	To        string
	Amount    int
	Signature string
}

/*
函数式选项模式只适合于复杂的实例化，如果参数只有简单几个，建议还是用普通的构造函数来解决。
*/
type TxOptions func(t *Tx)

func NewTx(from string, to string, amount int, signature string) *Tx {
	return &Tx{
		From:      from,
		To:        to,
		Amount:    amount,
		Signature: signature,
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

func WithSignature(signature string) TxOptions {
	return func(t *Tx) {
		t.Signature = signature
	}
}

func NewTxWithOptions(options ...TxOptions) *Tx {
	t := &Tx{}
	for _, option := range options {
		option(t)
	}
	return t
}

func (t *Tx) VerifyTx() bool {
	if strings.Contains(t.Signature, t.From) {
		return true
	}
	return false
}
