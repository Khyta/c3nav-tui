package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var apiEndpoint string = "https://39c3.c3nav.de/api/v2/"

// TODO: attach method to struct to store session key data
type sessionKey struct {
	Key string `json:"key"`
}

func (k *sessionKey) getSessionKey() {
	sessionURL := apiEndpoint + "auth/session/"
	resp, err := http.Get(sessionURL)
	if err != nil {
		slog.Error("response broken", "error", err)
		return
	}
	slog.Info("initial response.", "status", resp.Status)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response.", "error", err)
		return
	}

	err = json.Unmarshal(body, &k)
	if err != nil {
		slog.Error("could not unmarshal body json.", "error", err)
	}
}

func main() {

	var s sessionKey
	s.getSessionKey()
	slog.Info("got session key.", "key", s.Key)

	statusAPI := apiEndpoint + "auth/status"

	resp, err := http.Get(statusAPI)
	if err != nil {
		slog.Error("status check unsuccessfull.", "error", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body[:]))
}
