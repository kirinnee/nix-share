package main

import (
	"encoding/json"
	"os"
)

func getCacheValue(cachePath string) (map[string]StorePacketData, error) {
	body, err := os.ReadFile(cachePath)
	content := ""
	if err != nil {
		content = "{}"
	} else {
		content = string(body)
	}

	var ret map[string]StorePacketData
	e := json.Unmarshal([]byte(content), &ret)

	return ret, e
}

func setCacheValue(cachePath string, store map[string]StorePacketData) error {
	b, e := json.Marshal(store)
	if e != nil {
		return e
	}
	err := os.WriteFile(cachePath, b, 0644)
	return err
}

type StorePacketData struct {
	Host      string `json:"host"`
	PubKey    string `json:"pub_key"`
	TimeStamp int64  `json:"time_stamp"`
}
