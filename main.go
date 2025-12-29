package main

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

var apiEndpoint string = "https://39c3.c3nav.de/api/v2/"

type SessionKey struct {
	Key string `json:"key"`
}

type ApiStatus struct {
	KeyType  string   `json:"key_type"`
	Readonly bool     `json:"readonly"`
	Scopes   []string `json:"scopes"`
}

func (key *SessionKey) Fetch() error {
	sessionURL := apiEndpoint + "auth/session/"
	resp, err := http.Get(sessionURL)
	if err != nil {
		slog.Error("response broken", "error", err)
		return err
	}
	slog.Info("initial response.", "status", resp.Status)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response.", "error", err)
		return err
	}

	err = json.Unmarshal(body, &key)
	if err != nil {
		slog.Error("could not unmarshal body json.", "error", err)
		return err
	}
	return nil
}

func (status *ApiStatus) Check(key string) error {
	statusAPI := apiEndpoint + "auth/status"
	client := &http.Client{}

	req, err := http.NewRequest("GET", statusAPI, nil)
	if err != nil {
		slog.Error("could not form new request", "error", err)
	}

	// Critical Header for Auth recognition with the session key
	req.Header.Add("X-API-Key", key)
	resp, err := client.Do(req)

	// Need to check explicitly for status code as err is only for ISO Layers
	// 1-6 (not 7)
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		slog.Error("API didn't return an expected response.", "statuscode", resp.StatusCode, "key", key)
		err := "unreachable authentication status check. " + resp.Status
		return errors.New(err)
	} else if resp.StatusCode == 401 {
		slog.Error("not authorized to access API", "statuscode", resp.StatusCode, "key", key)
		err := "cannot access API for " + statusAPI + ". Got " + resp.Status
		return errors.New(err)
	}

	if err != nil {
		slog.Error("could not complete request with header.", "error", err)
		return errors.New("could not complete request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response.", "error", err)
		return errors.New("could not read response")
	}

	err = json.Unmarshal(body, &status)
	if err != nil {
		slog.Error("could not unmarshal body json response.", "error", err)
		return errors.New("could not unmarshal body json")
	}
	return nil
}

func main() {

	var session SessionKey
	err := session.Fetch()
	if err != nil {
		slog.Error("unable to get session key.")
		return
	}
	slog.Info("got session key.", "key", session.Key)

	var status ApiStatus
	err = status.Check(session.Key)
	if err != nil {
		slog.Error("could not get API status.", "error", err)
		return
	}
	slog.Info("got API status.", "status", status.KeyType)
}
