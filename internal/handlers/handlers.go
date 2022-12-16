package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/hoanguyen1998/crypto-payment-system/internal/blockchain"
	"github.com/hoanguyen1998/crypto-payment-system/internal/services"
)

type ServerHandler struct {
	services   *services.AppService
	redis      *redis.Client
	ethHandler *blockchain.EthHandler
}

func NewServerHandler(services *services.AppService, redis *redis.Client) *ServerHandler {
	return &ServerHandler{
		services:   services,
		redis:      redis,
		ethHandler: blockchain.NewEthHandler(),
	}
}
