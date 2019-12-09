package cryptlib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
)

func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)
	return cipherText, nil
}

func AesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]
	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func Encrypt(rawData, key []byte) (string, error) {
	data, err := AesCBCEncrypt(rawData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func Decrypt(rawData string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := AesCBCDecrypt(data, key)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func EncryptStream(reader io.Reader, key []byte) (io.Reader, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	cipherText, err := Encrypt(data, key)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader([]byte(cipherText)), nil
}

func DecryptStream(reader io.Reader, key []byte) (io.Reader, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	cipherText, err := Decrypt(string(data), key)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader([]byte(cipherText)), nil
}

func EncryptFile(src, dst string, key []byte) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	cipherText, err := Encrypt(data, key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, []byte(cipherText), 0666)
}

func DecryptFile(src, dst string, key []byte) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	cipherText, err := Decrypt(string(data), key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, []byte(cipherText), 0666)
}
