package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hoanguyen1998/crypto-payment-system/helpers"
)

type NewMasterPublicKeyPayload struct {
	MasterPulicKey  string `json:"master_public_key"`
	UserId          int    `json:"user_id"`
	PaymentMethodId int    `json:"payment_method_id"`
}

func (s *ServerHandler) CreateMasterPublicKey(w http.ResponseWriter, r *http.Request) {
	var masterKeyPayload NewMasterPublicKeyPayload

	err := json.NewDecoder(r.Body).Decode(&masterKeyPayload)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, helpers.NewInternalServerError(err.Error()))
		return
	}

	masterKey, errRest := s.services.NewMasterPublicKey(masterKeyPayload.MasterPulicKey, masterKeyPayload.UserId, masterKeyPayload.PaymentMethodId)

	if errRest != nil {
		s.writeJSON(w, errRest.Status, errRest)
		return
	}

	s.writeJSON(w, http.StatusOK, masterKey)
}
