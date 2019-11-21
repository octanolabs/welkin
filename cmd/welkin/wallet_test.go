package main

import (
	"strings"
	"testing"
)

const (
	mnemonic = "walnut tower salon print rural upper all always year video alarm empty"
)

func TestGenerateMnemonic(t *testing.T) {
	tMnemonic := GenerateMnemonic()
	arr := strings.Split(tMnemonic, " ")
	len := len(arr)
	if len != 24 {
		t.Fatal("invalid number of words in mnemonic")
	}
	return
}

func TestCreateWallet(t *testing.T) {
	wallet := CreateWallet(mnemonic)
	address := GetAddressFromPath(wallet, "0")
	if address != "0xf176F44A36A1F81f8aD2723Db64b8751F6235D2d" {
		t.Fatal("0xf176F44A36A1F81f8aD2723Db64b8751F6235D2d not found at position 0")
	}
	return
}

func TestGetAddressFromPath(t *testing.T) {
	wallet := CreateWallet(mnemonic)
	address := GetAddressFromPath(wallet, "0")
	if address != "0xf176F44A36A1F81f8aD2723Db64b8751F6235D2d" {
		t.Fatal("0xf176F44A36A1F81f8aD2723Db64b8751F6235D2d not found at position 0")
	}
	address = GetAddressFromPath(wallet, "1")
	if address != "0xDd2717dfe605a52EeD8eDe92Ae5D44AbD2c7eD36" {
		t.Fatal("0xDd2717dfe605a52EeD8eDe92Ae5D44AbD2c7eD36 not found at position 1")
	}
	return
}
