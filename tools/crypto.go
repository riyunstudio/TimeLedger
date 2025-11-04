package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// 加解密相關工具

// sha256 加密
func (tl *Tools) Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

// md5 加密
func (tl *Tools) Md5(str string) string {
	hash := md5.Sum([]byte(str))
	// 哈希轉十六進制字串
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

// AES-128 CBC 加密
func (tl *Tools) AesEncryptCBC(origData []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("加密金鑰的長度必須為16,24,32, Err: " + err.Error())
	}

	origData = tl.pkcs5Padding(origData, block.BlockSize())
	encrypted := make([]byte, len(origData))

	// 如果没有给定 IV，则使用密钥的前 16 字节作为默认 IV
	if len(iv) == 0 {
		iv = key[:block.BlockSize()]
	}
	// 检查 IV 长度
	if len(iv) != block.BlockSize() {
		return "", errors.New("IV長度必須等於塊大小16字節")
	}
	// 创建 CBC 加密器并加密
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(encrypted, origData)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// AES-128 EBC 加密
func (tl *Tools) AesEncryptEBC(origData []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("加密金鑰的長度必須為16,24,32, Err: " + err.Error())
	}

	blockSize := block.BlockSize()
	origData = tl.pkcs5Padding(origData, blockSize)
	encrypted := make([]byte, len(origData))

	for bs, be := 0, blockSize; bs < len(origData); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(encrypted[bs:be], origData[bs:be])
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// AES-128 CBC 解密
func (tl *Tools) AesDecryptCBC(cipherText string, key []byte, iv []byte) ([]byte, error) {
	// 檢查輸入參數
	if cipherText == "" {
		return nil, errors.New("密文不能為空")
	}

	encrypted, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("加密金鑰的長度必須為16,24,32, Err: " + err.Error())
	}

	if len(iv) == 0 {
		iv = key[:block.BlockSize()]
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)

	return tl.pkcs5UnPadding(decrypted), nil
}

// AES-128 EBC 解密
func (tl *Tools) AesDecryptEBC(cipherText string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("加密金鑰的長度必須為16,24,32, Err: " + err.Error())
	}

	blockSize := block.BlockSize()
	encrypted, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	if len(encrypted)%blockSize != 0 {
		return nil, fmt.Errorf("解密錯誤, Err: cipherText is not a multiple of the block size")
	}

	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, blockSize; bs < len(encrypted); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	return tl.pkcs5UnPadding(decrypted), nil
}

func (tl *Tools) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (tl *Tools) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 增加輸入驗證，避免panic
	if length == 0 {
		return origData
	}

	unpadding := int(origData[length-1])
	// 驗證unpadding值的有效性
	if unpadding > length {
		return origData
	}

	return origData[:(length - unpadding)]
}
