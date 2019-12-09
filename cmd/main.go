package main

import (
	"flag"
	"fmt"

	"github.com/echoturing/crypt/cryptlib"
)

var (
	src    string
	dst    string
	key    string
	action string
)

type Action string

const (
	ActionUndefined Action = ""
	ActionEncrypt   Action = "encrypt"
	ActionDecrypt   Action = "decrypt"
)

func parseFlag() {
	flag.StringVar(&src, "src", "", "src file name")
	flag.StringVar(&dst, "dst", "", "dst file name")
	flag.StringVar(&key, "key", "", "the key")
	flag.StringVar(&action, "action", "", "action:(encrypt|decrypt)")
	flag.Parse()

}

func EncryptFile(src, dst string, kb []byte) error {
	fmt.Printf("encrypt %s to %s\n", src, dst)
	return cryptlib.EncryptFile(src, dst, kb)
}
func DecryptFile(src, dst string, kb []byte) error {
	fmt.Printf("decrypt %s to %s\n", src, dst)
	return cryptlib.DecryptFile(src, dst, kb)
}
func main() {
	parseFlag()
	var err error
	kb := []byte(key)
	switch Action(action) {
	default:
		err = EncryptFile(src, dst, kb)
	case ActionEncrypt:
		err = EncryptFile(src, dst, kb)
	case ActionDecrypt:
		err = DecryptFile(src, dst, kb)
	}
	if err != nil {
		panic(err)
	}
}
