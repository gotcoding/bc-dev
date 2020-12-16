package ac

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/gotcoding/blockchain/01_crypto/03_ac/secp256k1"
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
func TestBTCKeyPair(t *testing.T) {
	pub, priv := secp256k1.GenerateKeyPair()
	fmt.Printf("私钥: %x \n", priv)
	fmt.Printf("公钥: %x \n", pub)
}

func TestEtheremSig(t *testing.T) {
	str := "以太坊签名练习"
	// 1.对签名的数据进行哈希运算
	hashInstance := sha256.New()
	hashInstance.Write([]byte(str))
	hashed := hashInstance.Sum(nil)
	pub, priv := secp256k1.GenerateKeyPair()
	// 2.执行以太坊中的secp256k1签名运算
	sigBytes, err := secp256k1.Sign(hashed, priv)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("签名为：%x ，长度为：%d \n", sigBytes, len(sigBytes))

	// 3.验证签名
	success := secp256k1.VerifySignature(pub, hashed, sigBytes[:len(sigBytes)-1])
	if !success {
		t.Error("验证签名失败！")
	} else {
		fmt.Println("验证签名成功！")
	}
	// 4.另外一种验证方式：恢复公钥后进行验证
	pubKey, err := secp256k1.RecoverPubkey(hashed, sigBytes)
	if err != nil {
		t.Error(err)
	}
	if !equal(pub, pubKey) {
		t.Error("验证签名失败！")
	} else {
		fmt.Println("验证签名成功！")
	}
}

func TestBitcoinSig(t *testing.T) {
	str := "比特币签名练习"
	// 1.对签名的数据进行哈希运算
	hashInstance := sha256.New()
	hashInstance.Write([]byte(str))
	hashed := hashInstance.Sum(nil)
	pub, priv := secp256k1.GenerateKeyPair()
	// 2.执行以太坊中的secp256k1签名运算
	sigBytes, err := secp256k1.BitcoinSign(hashed, priv)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("签名为：%x ，长度为：%d \n", sigBytes, len(sigBytes))

	// 3.验证签名
	flag := secp256k1.BitcoinVerify(pub, hashed, sigBytes)
	fmt.Println("验证签名结果：", flag)
}

func equal(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
