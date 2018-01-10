package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"errors"
)

type MetaData struct {
	Sid    string `json:"sid"`
	Skey   string `json:"skey"`
	Region string `json:"region"`
}

type Cluster struct {
	ID string `json:"clusterid"`
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

func getClusterID() (err error) {
	if usr.Region == "" {
		err = errors.New("Choose Region First")
		return
	}
	path := getAPI("/cloud/cluster/info")
	data, err := SandHttp(http.MethodGet, path+"?region="+usr.Region)
	if err != nil {
		fmt.Printf("[Get Cluster Info Error] [%s]\n", err.Error())
		return
	}

	var cs []Cluster

	err = json.Unmarshal(data, &cs)
	if err != nil {
		fmt.Printf("[Parse Cluster Info Error] [%s]\n", err.Error())
		return
	}

	if len(cs) == 1 {
		usr.ClusterID = cs[0].ID
	} else {
		var css []string
		for _, s := range cs {
			css = append(css, s.ID)
		}

		usr.ClusterID, err = prompt(CLULabel, css)
		if err != nil {
			fmt.Printf("[Choose Cluster Info Error] [%s]\n", err.Error())
			return
		}

	}
	return
}
