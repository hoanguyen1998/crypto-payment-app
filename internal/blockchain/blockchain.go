package blockchain

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateAddress(paymentMethodId int, account string, index int) (string, error) {
	newIndex := uint32(index)

	if paymentMethodId == 2 {
		return generateEthereumAddress(account, newIndex)
	}
	return generateBtcAddress(account, newIndex)
}

func GenerateApplicationKey(paymentMethodId int, masterPubKeyStr string, index int) (string, error) {
	newIndex := uint32(index)
	return generateAccountFromMasterPubKey(masterPubKeyStr, newIndex)
}

func generateBtcAddress(account string, index uint32) (string, error) {
	accountKey, err := hdkeychain.NewKeyFromString(account)
	if err != nil {
		return "", err
	}
	childKey, _ := accountKey.Derive(index)
	address, _ := childKey.Address(&chaincfg.MainNetParams)
	return address.String(), nil
}

func generateAccountFromMasterPubKey(masterPubKeyStr string, index uint32) (string, error) {
	masterPubKey, err := hdkeychain.NewKeyFromString(masterPubKeyStr)
	if err != nil {
		return "", err
	}
	account, err := masterPubKey.Derive(index)

	if err != nil {
		return "", err
	}

	return account.String(), nil
}

func generateEthereumAddress(account string, index uint32) (string, error) {
	accountKey, err := hdkeychain.NewKeyFromString(account)
	if err != nil {
		return "", err
	}
	childKey, _ := accountKey.Derive(index)

	ecPubKey, _ := childKey.ECPubKey()

	pubBytes := ecPubKey.SerializeUncompressed()

	publicKey, _ := crypto.UnmarshalPubkey(pubBytes)

	ethAddress := crypto.PubkeyToAddress(*publicKey).String()

	return ethAddress, nil
}
