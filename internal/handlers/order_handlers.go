package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hoanguyen1998/crypto-payment-system/helpers"
)

type NewOrderPayload struct {
	ApplicationOrderId int `json:"application_order_id"`
	Amount             int `json:"amount"`
	PaymentMethodId    int `json:"payment_method_id"`
	AppId              int `json:"app_id"`
}

func (s *ServerHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()
	var orderPayload NewOrderPayload

	err := json.NewDecoder(r.Body).Decode(&orderPayload)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, helpers.NewInternalServerError(err.Error()))
		return
	}

	order, errRest := s.services.CreateOrder(orderPayload.AppId, orderPayload.PaymentMethodId, orderPayload.PaymentMethodId, float64(orderPayload.Amount))

	if errRest != nil {
		s.writeJSON(w, errRest.Status, errRest)
		return
	}

	s.redis.Set(ctx, order.ReceivedAddress, helpers.EtherToWei(order.Amount)+"#"+strconv.Itoa(order.Id), 0)

	s.writeJSON(w, http.StatusOK, order)
}

func (s *ServerHandler) GetOrdersToWithdraw(w http.ResponseWriter, r *http.Request) {
	keyIdStr := r.URL.Query().Get("appKeyId")

	appKeyId, err := strconv.Atoi(keyIdStr)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, helpers.NewInternalServerError(err.Error()))
		return
	}

	orders, errRest := s.services.GetOrdersToWithdraw(appKeyId)

	if errRest != nil {
		s.writeJSON(w, errRest.Status, errRest)
		return
	}

	s.writeJSON(w, http.StatusOK, orders)
}
