package main

import (
	"strings"
	"net/http"
	"fmt"

	"github.com/urfave/cli"
	"os"
	"github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"encoding/json"
)

func nsAction(c *cli.Context) error {
	getSIP(c)
	getOpt := c.String("get")
	if getOpt != "" {
		getAllNamespace(getOpt)
	}

	return nil
}

func getAllNamespace(ns string) {
	//label := "Choose The Namespace"
	var nssn []string

	if strings.ToLower(ns) == "all" {
		data, err := SandHttp(http.MethodGet, getAPI("/cloud/namespace/info")+"?clusterid=cls-rfje0azd")
		if err != nil {
			fmt.Printf("[%s]\n", err)
			os.Exit(-1)
		}
		var nss []namespace.NSInfo_data_namespaces

		err = json.Unmarshal(data, &nss)

		for _, n := range nss {
			nssn = append(nssn, n.Name)
		}

	} else {
		data, err := SandHttp(http.MethodGet, getAPI("/cloud/namespace/info")+"?clusterid=cls-rfje0azd&name="+ns)
		if err != nil {
			os.Exit(-1)
		}
		var nss namespace.NSInfo_data_namespaces

		err = json.Unmarshal(data, &nss)
		nssn = append(nssn, nss.Name)
	}

	cns, err := prompt(NSLabel, nssn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cns)
}
