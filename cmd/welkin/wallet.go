package main

import (
	"log"

	"github.com/miguelmota/go-ethereum-hdwallet"
)

/*GenerateMnemonic ...
params:
returns: string, 24 word mnemonic phrase.*/
func GenerateMnemonic() string {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		return ""
	}

	return mnemonic
}

/*CreateWallet ...
params: string, 24 word mnemonic phrase.
returns: Wallet, https://godoc.org/github.com/miguelmota/go-ethereum-hdwallet#Wallet*/
func CreateWallet(mnemonic string) *hdwallet.Wallet {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
	return wallet
}

/*GetAddressFromPath ...
params: Wallet.
returns: Account, https://godoc.org/github.com/ethereum/go-ethereum/accounts#Account*/
func GetAddressFromPath(wallet *hdwallet.Wallet, uid string) string {
	path := hdwallet.MustParseDerivationPath("m/44'/108'/0'/0/" + uid)
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	return account.Address.Hex()
}
