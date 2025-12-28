package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var apiEndpoint string = "https://38c3.c3nav.de/api/v2/"

func main() {

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
