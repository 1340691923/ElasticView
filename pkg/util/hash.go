package util

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Hash 散列函数 返回空值则为错误
func Hash(ht crypto.Hash, b []byte) []byte {
	switch ht {
	case crypto.MD5:
		h := md5.New()
		h.Write(b)
		return h.Sum(nil)
	case crypto.SHA1:
	case crypto.SHA256:
	}
	return nil
}

// HashHex can hash and encode to string
func HashHex(ht crypto.Hash, b []byte) string {
	return hex.EncodeToString(Hash(ht, b))
}

// MD5Hash md5
func MD5Hash(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

// MD5HexHash md5 and encode to string
func MD5HexHash(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256HexHash sha256 and encode to string
func SHA256HexHash(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1HexHash sha1 and encode to string
func SHA1HexHash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA1Hash hmac-sha1
func HMACSHA1Hash(src, key []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(src)
	return h.Sum(nil)
}

// HMACSHA1HexHash hmac-sha1 and encode to string
func HMACSHA1HexHash(src, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func PwdEncode(pwd string, key string) string {
	strbytes := []byte(fmt.Sprintf("%s:%s", pwd, key))
	encoded := base64.StdEncoding.EncodeToString(strbytes)
	return encoded
}

func PwdDecode(pwd string, key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%s:%s", pwd, key))
	if err != nil {
		return "", nil
	}
	return string(decoded), nil
}

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}

func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
