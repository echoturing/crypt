package modules

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func GenerateKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publickey := &privatekey.PublicKey
	return privatekey, publickey, nil
}

func Encrypt(plaintext string, publickey *rsa.PublicKey) (string, error) {
	label := []byte("")
	sha256hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(sha256hash, rand.Reader, publickey, []byte(plaintext), label)

	decodedtext := base64.StdEncoding.EncodeToString(ciphertext)
	return decodedtext, err
}
func Decrypt(ciphertext string, privatekey *rsa.PrivateKey) (string, error) {
	decodedtext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	sha256hash := sha256.New()
	decryptedtext, err := rsa.DecryptOAEP(sha256hash, rand.Reader, privatekey, decodedtext, nil)
	if err != nil {
		return "", fmt.Errorf("RSA decrypt failed, error=%s\n", err.Error())
	}

	return string(decryptedtext), nil
}

func DumpPrivateKeyFile(privatekey *rsa.PrivateKey, filename string) error {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keybytes,
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func DumpPublicKeyFile(publickey *rsa.PublicKey, filename string) error {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keybytes,
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
func LoadPrivateKeyFile(keyfile string) (*rsa.PrivateKey, error) {
	keybuffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(keybuffer))
	if block == nil {
		return nil, errors.New("private key error!")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}

func LoadPublicKeyFile(keyfile string) (*rsa.PublicKey, error) {
	keybuffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keybuffer)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}

func DumpPrivateKeyBuffer(privatekey *rsa.PrivateKey) (string, error) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keybytes,
	}

	var keybuffer []byte = pem.EncodeToMemory(block)
	return string(keybuffer), nil
}
func DumpPublicKeyBuffer(publickey *rsa.PublicKey) (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keybytes,
	}

	var keybuffer []byte = pem.EncodeToMemory(block)
	return string(keybuffer), nil
}

func LoadPrivateKeyBase64(base64key string) (*rsa.PrivateKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}
func LoadPublicKeyBase64(base64key string) (*rsa.PublicKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(keybytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}
