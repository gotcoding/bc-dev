package ac

import (
	"fmt"
	"testing"
)

func TestGenerateRSA(t *testing.T) {
	if err := GenerateRSAKey(); err != nil {
		t.Error(err)
		t.Fatal("密钥生成失效")
	}
	t.Log("密钥生成成功！")
}

func TestRSAEncryptAndDecrypt(t *testing.T) {
	str := "RSA密码练习"
	fmt.Println("原始字符串：", str)
	cipherBytes, err := RSAEncrypt([]byte(str), "./files/public.pem")
	if err != nil {
		t.Fatal(err)
	}
	originalBytes, err := RSADecrypt(cipherBytes, "./files/private.pem")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("解密后的字符串：", string(originalBytes))
}

func TestRSAString(t *testing.T) {
	str := "RSA密码练习"
	fmt.Println("原始字符串：", str)
	cipherText, err := RSAEncryptString(str, "./files/public.pem")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("加密后的字符串：", cipherText)
	originalText, err := RSADecryptString(cipherText, "./files/private.pem")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("解密后的字符串：", originalText)
}

func TestRASSign(t *testing.T) {
	str := "RSA签名"
	base64Sig, _ := RSASign([]byte(str), "./files/private.pem")
	fmt.Println("签名后信息:", base64Sig)
	err := RSAVerify([]byte(str), base64Sig, "./files/public.pem")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("签名验证成功！")
}
