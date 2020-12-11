package hash_dev

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	str := "MD5测试"
	md5_1(str)
	md5_2(str)
	md5_3(str)
}

func TestSum(t *testing.T) {
	str := "MD5测试"
	pre := "test"
	sum(str, pre)
}

func TestStringHash(t *testing.T) {
	b1, _ := hex.DecodeString("6da0d1d5")
	b2, _ := hex.DecodeString("6da0d155")
	fmt.Printf("%x\n", md5.Sum(b1))
	fmt.Printf("%x\n", md5.Sum(b2))
	b3 := []byte("6da0d1d5")
	b4 := []byte("6da0d155")
	fmt.Printf("%x\n", md5.Sum(b3))
	fmt.Printf("%x\n", md5.Sum(b4))
}
