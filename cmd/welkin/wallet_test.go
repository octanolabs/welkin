package main

import (
	"strings"
	"testing"
)

func TestGenerateMnemonic(t *testing.T) {
	mnemonic := GenerateMnemonic()
	arr := strings.Split(mnemonic, " ")
	len := len(arr)
	if len != 24 {
		t.Fatal("invalid number of words in mnemonic")
	}
	return
}

func TestCreateWallet(t *testing.T) {
	mnemonic := "walnut tower salon print rural upper all always year video alarm empty"
	wallet := CreateWallet(mnemonic)
	address := GetAddressFromPath(wallet, "0")
	if address != "0xf176F44A36A1F81f8aD2723Db64b8751F6235D2d" {
		t.Fatal("0xf176F44A36A1F81f8aD2723Db64b8751F6235D2d not found at position 0")
	}
	return
}
