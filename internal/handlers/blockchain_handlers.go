package handlers

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/hoanguyen1998/crypto-payment-system/internal/blockchain"
	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

type SendTransactionPayload struct {
	TxHash          string `json:"tx_hash"`
	PaymentMethodId int    `json:"payment_method_id"`
}

func (s *ServerHandler) SendTransaction(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()
	var txPayload SendTransactionPayload

	json.NewDecoder(r.Body).Decode(&txPayload)

	s.redis.Set(ctx, txPayload.TxHash, txPayload.PaymentMethodId, 0)
}

func (s *ServerHandler) ProcessEthBlock() {
	number, _ := s.services.GetLatestBlockNumber(1)

	var blockNumber *big.Int

	var txs []blockchain.EthTransaction

	if number != 0 {
		blockNumber = big.NewInt(int64(number))
	} else {
		blockNumber, _ = s.ethHandler.GetLatestBlock()
	}

	txs, _ = s.ethHandler.GetBlockByNumber(blockNumber)

	for {
		s.HandleEthTransaction(txs)

		blockNumber = new(big.Int).Add(blockNumber, big.NewInt(int64(1)))
		txs, _ = s.ethHandler.GetBlockByNumber(blockNumber)

		if txs == nil {
			break
		}
	}
}

func (s *ServerHandler) HandleEthTransaction(txs []blockchain.EthTransaction) {
	var ctx = context.Background()

	for _, tx := range txs {
		orderIdStr, _ := s.redis.Get(ctx, tx.Hash).Result()
		orderId, err := strconv.Atoi(orderIdStr)
		if err == nil {
			s.services.UpdateOrderStatusAndCreateTransaction(orderId, "withdraw", models.Transaction{TxHash: tx.Hash})
			continue
		}

		amountVal, _ := s.redis.Get(ctx, tx.To).Result()
		amountArr := strings.Split(amountVal, "#")
		amountStr := amountArr[0]
		amount, _ := new(big.Int).SetString(amountStr, 0)
		amountToCmp, _ := new(big.Int).SetString(tx.Value, 0)
		if amount.Cmp(amountToCmp) == 0 {
			orderId, err = strconv.Atoi(amountArr[1])
			if err == nil {
				s.services.UpdateOrderStatusAndCreateTransaction(orderId, "confirm", models.Transaction{TxHash: tx.Hash})
			}
		}
	}

	//https://blog.devgenius.io/big-int-in-go-handling-large-numbers-is-easy-157cb272dd4f
	// connvert normal eth to wei : num-of-eth * 10**18 (1 and 18 zeros)

	// amountStr, _ := s.redis.Get(ctx, tx.To).Result()
	// amount, _ := new(big.Int).SetString(amountStr, 0)
	// amountToCmp, _ := new(big.Int).SetString(tx.Value, 0)
	// if amount.Cmp(amountToCmp) == 0 {
	// 	orderIdStr, _ = s.redis.Get(ctx, tx.To+"#"+amountStr).Result()
	// 	orderId, err = strconv.Atoi(orderIdStr)
	// 	if err == nil {
	// 		s.services.UpdateOrderStatusAndCreateTransaction(orderId, "confirm", models.Transaction{TxHash: tx.Hash})
	// 	}
	// }
}
