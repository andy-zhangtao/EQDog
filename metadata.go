package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
)

type MetaData struct {
	Sid    string `json:"sid"`
	Skey   string `json:"skey"`
	Region string `json:"region"`
}

var md []MetaData

func syncMetaData() (err error) {
	data, err := SandHttp(http.MethodGet, getAPI("/cloud/metadata"))
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &md)
	var rn []string
	for _, m := range md {
		rn = append(rn, m.Region)
	}

	region, err := prompt(RGLabel, rn)
	if err != nil {
		fmt.Println(err)
	}

	usr.Region = region
	return
}

func isMdValid() {
	//if md.Sid == "" {
	//	syncMetaData()
	//}
}

func metaAction(c *cli.Context) error {
	getSIP(c)
	syncMetaData()
	return nil
}
