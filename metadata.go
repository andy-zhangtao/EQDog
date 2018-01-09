package main

import (
	"net/http"
	"encoding/json"
)

type MetaData struct {
	Sid    string `json:"secret_id"`
	Skey   string `json:"secret_key"`
	Region string `json:"region"`
}

var md MetaData

func syncMetaData() (err error) {
	data, err := SandHttp(http.MethodGet, getAPI("/cloud/metadata")+"?region=sh")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &md)

	return
}

func isMdValid() {
	if md.Sid == "" {
		syncMetaData()
	}
}
