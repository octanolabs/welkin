package main

import (
	"github.com/miguelmota/go-ethereum-hdwallet"
)

func GenerateMnemonic() string {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		return ""
	}
	return mnemonic
}
