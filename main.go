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
	if resp.StatusCode != 200 {
		slog.Error("unexpected status code.", "status", resp.StatusCode)
		errmsg := "unexpected status code. cannot fetch session auth cookie. " + resp.Status
		return errors.New(errmsg)
	}
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
	raw, err := apiRequestRaw(key, "auth/status")
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, status)
}

// func getLocation(key SessionKey) {
// 	locationsEndpoint := apiEndpoint + "/map/locations/"

// }

func apiRequestRaw(key string, specificEndpoint string) (json.RawMessage, error) {

	url := apiEndpoint + specificEndpoint
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("could not form new request", "error", err)
	}

	// Critical Header for Auth recognition with the session key
	req.Header.Add("X-API-Key", key)
	resp, err := client.Do(req)

	// Need to check explicitly for status code as err is only for ISO Layers
	// 1-6 (not 7)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			slog.Error("not authorized to access API", "statuscode", resp.StatusCode, "key", key)
			err := "cannot access API for " + url + ". Got " + resp.Status
			return nil, errors.New(err)
		}
		slog.Error("API didn't return an expected response.", "statuscode", resp.StatusCode, "key", key)
		err := "unreachable authentication status check. " + resp.Status
		return nil, errors.New(err)
	}

	if err != nil {
		slog.Error("could not complete request with header.", "error", err)
		return nil, errors.New("could not complete request")
	}
	defer resp.Body.Close()

	var rawMessage json.RawMessage
	err = json.NewDecoder(resp.Body).Decode(&rawMessage)
	if err != nil {
		slog.Error("could not decode raw message", "error", err)
		return nil, errors.New("could not read response")
	}

	return rawMessage, nil
}

func main() {

	var session SessionKey
	err := session.Fetch()
	if err != nil {
		slog.Error("unable to get session key.", "error", err)
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

	// TODO: Get locations (slim)
	//	- Parse locations
	// 	- Provide location description
	// TODO: Get regular Updates for auth with tileserver
}
