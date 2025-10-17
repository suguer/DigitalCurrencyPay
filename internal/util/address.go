package util

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func HexString2Address(hex string) (string, error) {
	if len(hex) == 40 {
		hex = "0x41" + hex
	} else if len(hex) == 42 && hex[:2] == "0x" {
		hex = hex[:2] + "41" + hex[2:]
	}
	b, err := HexStringToBytes(hex)
	if err != nil {
		return "", err
	}
	a := EncodeCheck(b)
	return a, nil
}

func Generate() (address string, key string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("无法生成私钥: %v", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	walletAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	address = walletAddress.Hex()
	key = hex.EncodeToString(privateKeyBytes)
	return address, key, nil
}

func MatchAddress(addr1, addr2 string) bool {
	addr1 = strings.ToUpper(addr1)
	addr2 = strings.ToUpper(addr2)
	if len(addr1) == 44 && strings.HasPrefix(addr1, "0x41") {
		addr1 = addr1[4:]
	}
	if len(addr2) == 44 && strings.HasPrefix(addr2, "0x41") {
		addr2 = addr2[4:]
	}
	return addr1 == addr2
}
