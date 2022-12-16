package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type EthHandler struct {
	EthClient *ethclient.Client
}

// Transaction data structure
type EthTransaction struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
	To    string `json:"to"`
}

func NewEthHandler() *EthHandler {
	client, err := ethclient.Dial("http://localhost:7545")

	if err != nil {
		fmt.Println(err)
	}

	return &EthHandler{EthClient: client}
}

func (e *EthHandler) GetBlockByNumber(blockNumber *big.Int) ([]EthTransaction, error) {
	block, err := e.EthClient.BlockByNumber(context.Background(), blockNumber)

	if err != nil {
		return nil, err
	}

	return e.HandleBlock(block), nil
}

func (e *EthHandler) GetLatestBlock() (*big.Int, error) {
	var ctx = context.Background()
	header, err := e.EthClient.HeaderByNumber(ctx, nil)

	if err != nil {
		return nil, err
	}

	return header.Number, nil
}

func (e *EthHandler) HandleBlock(block *types.Block) []EthTransaction {
	// We add a recover function from panics to prevent our API from crashing due to an unexpected error
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var ethTx []EthTransaction

	for _, tx := range block.Transactions() {
		ethTx = append(ethTx, EthTransaction{
			Hash:  tx.Hash().String(),
			Value: tx.Value().String(),
			To:    tx.To().String(),
		})
	}

	return ethTx
}

func SendEthRawTransaction(rawTx string) {
	client, err := ethclient.Dial("http://localhost:7545")

	if err != nil {
		fmt.Println(err)
	}

	rawTxBytes, err := hex.DecodeString(rawTx)

	if err != nil {
		fmt.Println(err)
	}

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())

	client.Close()
}
