package wallet

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"gochain/cryptor"
)

type Wallet struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	address    string
}

func (w *Wallet) GetAddress() string {
	return w.address
}

// NewWallet 创建一个新的钱包，生成密钥对
func NewWallet() *Wallet {
	privateKey, publicKey, err := cryptor.GenerateRsaKeyPair(2048)
	if err != nil {
		return nil
	}

	// 生成钱包地址，通常是公钥的哈希
	address := generateAddress(publicKey)

	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
		address:    address,
	}
}

// GetPublicKeyAsString 获取RSA公钥的字符串形式 (PEM格式)
func (w *Wallet) GetPublicKeyAsString() string {
	return cryptor.EncodePublicKey(w.publicKey)
}

// Sign 对交易数据进行签名
func (w *Wallet) Sign(data string) (string, error) {
	signature, err := cryptor.Sign(w.privateKey, data)
	if err != nil {
		return "", err
	}
	return signature, nil
}

// Verify 验证交易签名
func (w *Wallet) Verify(data, signature string) bool {
	valid, err := cryptor.VerifySig(w.publicKey, data, signature)
	if err != nil {
		log.Println("Verification failed:", err)
		return false
	}
	return valid
}

// generateAddress 通过公钥生成钱包地址
func generateAddress(publicKey *rsa.PublicKey) string {
	// X.509格式序列化公钥
	pubBytes := x509.MarshalPKCS1PublicKey(publicKey)

	// 对公钥进行SHA-256哈希
	hash := sha256.Sum256(pubBytes)

	// 将哈希值编码为Base58
	address := hex.EncodeToString(hash[:])
	return address
}

func (w *Wallet) SaveToFile(filename string) error {
	// 创建/打开文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 编码并写入私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(w.privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(file, privateKeyPEM)
	if err != nil {
		return err
	}

	return nil
}

func LoadWalletFromFile(filename string) (*Wallet, error) {
	// 读取文件内容
	keyPairPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解码PEM块（私钥）
	block, _ := pem.Decode(keyPairPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("无效的私钥数据")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		address:    generateAddress(&privateKey.PublicKey),
	}
	return wallet, nil
}
