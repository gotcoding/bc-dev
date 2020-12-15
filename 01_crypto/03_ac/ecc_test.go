package ac

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestECCSign(t *testing.T) {
	data := "ECC签名测试"
	hashInstance := sha256.New()
	hashInstance.Write([]byte(data))
	hashed := hashInstance.Sum(nil)
	// 生成公私钥对
	privateKey, publicKeyBytes := NewKeyPair()
	// 生成DER格式签名
	derSign := ECDSASign(hashed, privateKey)
	fmt.Println("签名信息：", derSign)
	// 验证签名
	flag := ECDSAVerify(publicKeyBytes, hashed, derSign)
	fmt.Println("验证签名结果：", flag)
}

// 校验地址 https://gobittest.appspot.com/Address
func TestGenerateKeyPair(t *testing.T) {
	pub, priv := GenerateKeyPair()
	fmt.Printf("私钥: %x \n", priv)
	fmt.Printf("公钥: %x \n", pub)
}
