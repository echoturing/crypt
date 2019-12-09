package cryptlib

import "testing"

func TestCrypt(t *testing.T) {
	key := []byte("AVy8pVVX2HEyiucWnzBiwDhrqLx2gsbY")
	plainText := "this is a plain text!"
	cipherText, err := Encrypt([]byte(plainText), key)
	if err != nil {
		t.Error(err)
		return
	}

	decryptPlainText, err := Decrypt(cipherText, key)
	if err != nil {
		t.Error(err)
		return
	}
	if plainText != decryptPlainText {
		t.Error("aes test failed")
	}
}
