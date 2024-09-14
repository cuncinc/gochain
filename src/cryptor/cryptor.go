package cryptor

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// Sha256 计算字符串的 SHA-256 哈希值
func Sha256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash) // 转换为 16 进制字符串
}

// GenerateRsaKeyPair 生成 RSA 密钥对
func GenerateRsaKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// Sig 对消息进行签名并返回可读的 Base64 字符串
func Sign(privateKey *rsa.PrivateKey, msg string) (string, error) {
	hashed := sha256.Sum256([]byte(msg))
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return "", err
	}
	// 使用 Base64 编码签名
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySig 验证签名（Base64 格式）
func VerifySig(publicKey *rsa.PublicKey, msg string, signatureBase64 string) (bool, error) {
	hashed := sha256.Sum256([]byte(msg))
	// 先解码 Base64 签名
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}
	// 验证签名
	err = rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, nil)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %v", err)
	}
	return true, nil
}

// EncodePublicKey 将 RSA 公钥编码为 PEM 格式的 Base64 字符串
func EncodePublicKey(pub *rsa.PublicKey) string {
	// 先序列化 RSA 公钥
	pubBytes := x509.MarshalPKCS1PublicKey(pub)
	// 再转换为 PEM 格式
	pubPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	}
	// 最后使用 Base64 编码
	return base64.StdEncoding.EncodeToString(pem.EncodeToMemory(pubPEM))
}

// DecodePublicKey 将 PEM 格式的 Base64 字符串解码为 RSA 公钥
func DecodePublicKey(pubStr string) (*rsa.PublicKey, error) {
	// 先解码 Base64 公钥字符串
	pubBytes, err := base64.StdEncoding.DecodeString(pubStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}
	// 解码 PEM 格式
	block, _ := pem.Decode(pubBytes)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid public key data")
	}
	// 解析 RSA 公钥
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	return pub, nil
}

// VerifySignatureWithPublicKeyString 使用RSA公钥字符串验证签名
func VerifySignatureWithPublicKeyString(pubKeyStr string, msg string, signatureBase64 string) (bool, error) {
	pub, err := DecodePublicKey(pubKeyStr)
	if err != nil {
		return false, err
	}

	return VerifySig(pub, msg, signatureBase64)
}
