package main

import (
	"encoding/json"
	"errors"
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

func getApiStatus(key string) (string, error) {
	statusAPI := apiEndpoint + "auth/status"
	client := &http.Client{}

	req, err := http.NewRequest("GET", statusAPI, nil)
	if err != nil {
		slog.Error("could not form new request", "error", err)
	}

	req.Header.Add("X-API-Key", key)
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		slog.Error("API didn't return a 200 OK.", "statuscode", resp.StatusCode)
		err := "unreachable API. " + resp.Status
		return "", errors.New(err)
	}

	if err != nil {
		slog.Error("could not complete request with header.", "error", err)
		return "", errors.New("could not complete request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response.", "error", err)
		return "", errors.New("could not read response")
	}

	return string(body[:]), nil
}

func main() {

	var s sessionKey
	s.getSessionKey()
	slog.Info("got session key.", "key", s.Key)
	apiStatus, err := getApiStatus(s.Key)
	if err != nil {
		slog.Error("could not get API status.", "error", err)
		return
	}
	slog.Info("get API status.", "status", apiStatus)
}
