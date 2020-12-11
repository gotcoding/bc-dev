package hash_dev

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

/*
Go的MD5包是采用RFC1321算法,
MD5加密算法已经被破解，不要用于安全应用程序。
*/

// MD5的 3 种使用方式
// 方式一：直接计算，比较简洁。
// 缺点是如果长度超出范围，不会报错。
func md5_1(str string) {
	res := md5.Sum([]byte(str))
	fmt.Println(hex.EncodeToString(res[:]))
}

// 方式二：先初始化后逐步写入
func md5_2(str string) {
	md := md5.New()              // 初始化一个MD5对象
	_, _ = md.Write([]byte(str)) // 写入数据，会检查输入长度是否超出
	res := md.Sum(nil)           // 执行计算
	fmt.Printf("%x\n", res)
	fmt.Println(hex.EncodeToString(res))
}

// 方式三：适合文件写入的方式
func md5_3(str string) {
	md := md5.New()
	_, _ = io.WriteString(md, str)
	fmt.Println(hex.EncodeToString(md.Sum(nil)))
}

// 查看源码，Sum函数将传入的参数直接加在hash结果的前面形成一个新的byte切片。
// 因此通常的使用方法就是将data置为nil。
func sum(str string, pre string) {
	m1 := md5.New()              // 初始化一个MD5对象
	_, _ = m1.Write([]byte(str)) // 写入数据，会检查输入长度是否超出
	res := m1.Sum(nil)           // 执行计算
	fmt.Println(hex.EncodeToString(res))

	fmt.Printf("%x\n", pre)
	res2 := m1.Sum([]byte(pre))
	fmt.Println(hex.EncodeToString(res2))
}

// 判断输入数据是否为16进制的数据
func MD5(str string, isHex bool) string {
	hashInstance := md5.New()
	if isHex {
		arr, _ := hex.DecodeString(str)
		hashInstance.Write(arr)
	} else {
		hashInstance.Write([]byte(str))
	}
	res := hashInstance.Sum(nil)
	return fmt.Sprintf("%x", res)
}

// 两次MD5
func MD5Double(str string, isHex bool) string {
	hashInstance := md5.New()
	if isHex {
		arr, _ := hex.DecodeString(str)
		hashInstance.Write(arr)
	} else {
		hashInstance.Write([]byte(str))
	}
	res := hashInstance.Sum(nil)
	hashInstance.Reset() //重置
	hashInstance.Write(res)
	res = hashInstance.Sum(nil)
	return fmt.Sprintf("%x", res)
}
