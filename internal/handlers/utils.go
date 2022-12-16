package handlers

import (
	"encoding/json"
	"net/http"
)

func (app *ServerHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
