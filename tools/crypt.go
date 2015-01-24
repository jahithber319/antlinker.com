// crypt
package tools

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

var coder = base64.StdEncoding

// Base64编码
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64解码
func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

// 返回MD5加密串
func MD5(src []byte) string {
	h := md5.New()
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}
