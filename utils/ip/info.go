package ip

//package main

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	url = "https://ipinfo.io/"
)

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Timezone string `json:"timezone"`
}

func Info(ip string) (*IPInfo, error) {
	response, err := http.Get(url + ip)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data IPInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// func main() {
// 	Info("77.91.123.163")
// }
