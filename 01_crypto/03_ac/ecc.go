package ac

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

// 生成公私钥对，
func NewKeyPair() (*ecdsa.PrivateKey, []byte) {
	// 1.生成椭圆曲线对象
	curve := elliptic.P256()
	// 2.生成密钥对，返回私钥对象
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 3.编码生成公私钥字节数组
	publicBytes := elliptic.Marshal(curve, privateKey.X, privateKey.Y)
	fmt.Printf("公钥：%x\n", publicBytes)
	return privateKey, publicBytes
}

// ECDSA数字签名
func ECDSASign(hashed []byte, privateKey *ecdsa.PrivateKey) string {
	// 1.数字签名生成r,s的big.Int对象
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return ""
	}
	// 2.将r,s转成r,s字符串
	strSigR := fmt.Sprintf("%x", r)
	strSinS := fmt.Sprintf("%x", s)
	if len(strSigR) == 63 {
		strSigR = "0" + strSigR
	}
	if len(strSinS) == 63 {
		strSinS = "0" + strSinS
	}
	// 3.r和s字符串拼接，形成数字签名的der格式
	derString := MakeDERSignString(strSigR, strSinS)
	return derString
}

// 生成数字签名的DER格式
func MakeDERSignString(strR, strS string) string {
	// 使用比特币格式
	// 获取R和S的长度
	lenSinR := len(strR) / 2
	lenSinS := len(strS) / 2

	// 2.计算DER序列的总长度
	len := lenSinR + lenSinS + 4

	// 3.将10进制长度转16进制字符串
	strLenSinR := fmt.Sprintf("%x", int64(lenSinR))
	strLenSinS := fmt.Sprintf("%x", int64(lenSinS))
	strLen := fmt.Sprintf("%x", int64(len))

	// 4.拼凑DER编码格式
	derString := "30" + strLen
	derString += "02" + strLenSinR + strR
	derString += "02" + strLenSinS + strS
	derString += "01"
	return derString
}

// ECDSA验证签名
func ECDSAVerify(publicKeyBytes, hashed []byte, derSignString string) bool {
	// 公钥长度
	keyLen := len(publicKeyBytes)
	if keyLen != 65 {
		return false
	}
	// 1.生成椭圆曲线对象
	curve := elliptic.P256()
	// 2.根据公钥字节数字，获取公钥的x和y
	publicKeyBytes = publicKeyBytes[1:]
	x := new(big.Int).SetBytes(publicKeyBytes[:32])
	y := new(big.Int).SetBytes(publicKeyBytes[32:])
	// 3.生成公钥对象
	publicKey := &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	// 4.对der格式的签名进行解析，获取r和s字节数组后转成big.Int
	rBytes, sBytes := ParseDERSignString(derSignString)
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)
	// 5.验证签名
	return ecdsa.Verify(publicKey, hashed, r, s)
}

func ParseDERSignString(derString string) (rBytes, sBytes []byte) {
	derBytes, _ := hex.DecodeString(derString)
	rBytes = derBytes[4:36]
	sBytes = derBytes[len(derBytes)-33 : len(derBytes)-1]
	return
}
