package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hoanguyen1998/crypto-payment-system/helpers"
)

type NewAppPayload struct {
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}

type NewApplicationKeyPayload struct {
	AppId           int `json:"app_id"`
	PaymentMethodId int `json:"payment_method_id"`
}

func (s *ServerHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	var appPayload NewAppPayload

	err := json.NewDecoder(r.Body).Decode(&appPayload)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, helpers.NewInternalServerError(err.Error()))
		return
	}

	app, errRest := s.services.NewApp(appPayload.Name, appPayload.UserId)

	if errRest != nil {
		s.writeJSON(w, errRest.Status, errRest)
		return
	}

	s.writeJSON(w, http.StatusOK, app)
}

func (s *ServerHandler) CreateAppKey(w http.ResponseWriter, r *http.Request) {
	var keyPayload NewApplicationKeyPayload

	err := json.NewDecoder(r.Body).Decode(&keyPayload)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, helpers.NewInternalServerError(err.Error()))
		return
	}

	appKey, errRest := s.services.NewAppKey(keyPayload.AppId, keyPayload.PaymentMethodId)

	if errRest != nil {
		s.writeJSON(w, errRest.Status, errRest)
		return
	}

	s.writeJSON(w, http.StatusOK, appKey)
}
