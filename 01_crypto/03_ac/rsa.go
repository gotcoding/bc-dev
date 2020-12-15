package ac

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	"io/ioutil"
	"os"
)

func GenerateRSAKey() error {
	// 生成RSA密钥对
	var bits int
	flag.IntVar(&bits, "key flag", 1024, "密钥长度，默认值是1024位")
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// 将私钥对象转换为DER编码形式
	derPrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	// 创建私钥pem文件
	file, err := os.Create("./files/private.pem")
	if err != nil {
		return err
	}
	// 对私钥信息进行编码，写入私钥文件中
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateKey,
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	file, err = os.Create("./files/public.pem")
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPublicKey,
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// RSA 加密
func RSAEncrypt(originalBytes []byte, fileName string) ([]byte, error) {
	// 1.读取公钥文件，解析出公钥对象
	publicKey, err := ReadParsePublicKey(fileName)
	if err != nil {
		return nil, err
	}
	// 2.RSA加密
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, originalBytes)
}

// RSA 解密
func RSADecrypt(cipherBytes []byte, fileName string) ([]byte, error) {
	// 1.读取私钥文件，解析出私钥对象
	privateKey, err := ReadParsePrivateKey(fileName)
	if err != nil {
		return nil, err
	}
	// 2.RSA解密
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
}

// RSA加密字符串，返回Base64处理的字符串
func RSAEncryptString(originalBytes, fileName string) (string, error) {
	cipherBytes, err := RSAEncrypt([]byte(originalBytes), fileName)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// RSA解密，经过base64的加密字符串，返回加密前的明文
func RSADecryptString(cipherText, fileName string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	originalBytes, err := RSADecrypt(cipherBytes, fileName)
	if err != nil {
		return "", err
	}
	return string(originalBytes), nil
}

// 读取公钥文件，解析出公钥对象
func ReadParsePublicKey(fileName string) (*rsa.PublicKey, error) {
	// 1.读取密钥文件，获取公钥字节
	publicKeyBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// 2.解码公钥字节，生成加密块对象
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	// 3.解析OER编码的公钥，生成公钥接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 4.公钥接口转型成公钥对象
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}

// 读取私钥文件，解析出私钥对象
func ReadParsePrivateKey(fileName string) (*rsa.PrivateKey, error) {
	// 1.读取密钥文件，获取公钥字节
	privateKeyBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// 2.解码公钥字节，生成加密块对象
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	// 3.解析OER编码的公钥，生成公钥接口
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// RSA签名
func RSASign(data []byte, fileName string) (string, error) {
	// 1.选择hash算法，对要签名的数据进行hash运算
	myHash := crypto.SHA256
	hashInstance := myHash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)
	// 2.读取私钥文件，解析出私钥对象
	privateKey, err := ReadParsePrivateKey(fileName)
	if err != nil {
		return "", err
	}
	// 3.RSA数字签名
	bytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, myHash, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// RSA验证签名
func RSAVerify(data []byte, base64Sig, fileName string) error {
	// 1.对base64编码的签名内容进行解码，返回签名字节
	bytes, err := base64.StdEncoding.DecodeString(base64Sig)
	if err != nil {
		return err
	}
	// 2.选择hash算法，对需要签名的数据进行hash运算
	myHash := crypto.SHA256
	hashInstance := myHash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)
	// 3.读取公钥文件，解析出公钥对象
	publicKey, err := ReadParsePublicKey(fileName)
	if err != nil {
		return err
	}
	// 4.RSA验证数字签名
	return rsa.VerifyPKCS1v15(publicKey, myHash, hashed, bytes)
}
