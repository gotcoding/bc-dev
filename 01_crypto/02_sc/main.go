package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
)

// 对称加密
func main() {
	// DES密钥
	key := "12345678"
	// 3DES密钥
	// key = "123456781234567812345678"
	// AES密钥
	// key = "1234567812345678"
	keyBytes := []byte(key)
	str := "a"
	cipherArr, err := SCEncrypt([]byte(str), keyBytes, "des")
	if err != nil {
		panic(err)
	}
	fmt.Printf("加密后的字节数组: %v\n", cipherArr)
	fmt.Printf("加密后的16进制: %x\n", cipherArr)

	originalBytes, err := SCDecrypt(cipherArr, keyBytes, "des")
	if err != nil {
		panic(err)
	}
	fmt.Println("解密后：", string(originalBytes))

}

// 对称加密算法（symmetric cryptography）
func SCEncrypt(originalBytes, key []byte, scType string) ([]byte, error) {
	var err error
	var block cipher.Block
	// 1.实例化密码器block（参数为密钥）
	switch scType {
	case "des":
		block, err = des.NewCipher(key)
	case "3des":
		block, err = des.NewTripleDESCipher(key)
	case "aes":
		block, err = aes.NewCipher(key)
	}
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	fmt.Println("-----", blockSize)
	// 2.对明文填充字节（参数为原始字节切片和密码对象的区块个数）
	paddingBytes := PKCS5Padding(originalBytes, blockSize)
	fmt.Println("填充后字节切片：", paddingBytes)
	// 3.实例化加密模式（参数为密码对象和密钥）
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	fmt.Println("加密模式：", blockMode)
	// 4.对填充字节后的明文进行加密
	cipherBytes := make([]byte, len(paddingBytes))
	blockMode.CryptBlocks(cipherBytes, paddingBytes)
	return cipherBytes, nil
}

func SCEncryptString(originalText, key, scType string) (string, error) {
	cipherBytes, err := SCEncrypt([]byte(originalText), []byte(key), scType)
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(cipherBytes)
	return base64Str, err
}

func SCDencryptString(originalText, key, scType string) (string, error) {
	cipherBytes, _ := base64.StdEncoding.DecodeString(originalText)
	cipherBytes, err := SCDecrypt(cipherBytes, []byte(key), scType)
	if err != nil {
		return "", err
	}
	return string(cipherBytes), err
}

func SCDecrypt(cipherBytes, key []byte, scType string) ([]byte, error) {
	var err error
	var block cipher.Block
	// 1.实例化密码器block（参数为密钥）
	switch scType {
	case "des":
		block, err = des.NewCipher(key)
	case "3des":
		block, err = des.NewTripleDESCipher(key)
	case "aes":
		block, err = aes.NewCipher(key)
	}
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// 2、实例化解密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 3、对密文进行解密
	paddingBytes := make([]byte, len(cipherBytes))
	blockMode.CryptBlocks(paddingBytes, cipherBytes)
	// 4、去除填充的字节
	originalBytes := PKCS5Padding(paddingBytes, blockSize)
	return originalBytes, nil
}

// 填充字节函数：数字填充
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	s := []byte{byte(padding)}
	s = bytes.Repeat(s, padding)
	return append(data, s...)
}

// 填充字节函数：零填充
func ZerosPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	s := []byte{0}
	s = bytes.Repeat(s, padding)
	return append(data, s...)
}

// 去除填充字节函数
func PKCS5UnPadding(data []byte) []byte {
	unPadding := data[len(data)-1]
	return data[:len(data)-int(unPadding)]
}

func ZerosUnPadding(data []byte) []byte {
	return bytes.TrimRightFunc(data, func(r rune) bool {
		return r == 0
	})
}
