package main

import (
	"strings"
	"net/http"
	"fmt"

	"github.com/urfave/cli"
	"os"
	"github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"encoding/json"
	"github.com/manifoldco/promptui"
	gs "github.com/andy-zhangtao/gogather/strings"
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
	usr.Namespace = cns
}

func nsOperation() {
	for {
		prompt := promptui.Prompt{
			Label:     "EQCloud Namespace >",
			IsVimMode: true,
			AllowEdit: true,
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		if ret := switchNS(result); ret {
			return
		}
	}
}

func switchNS(cmd string) (ret bool) {
	ret = false

	final := gs.RemoveMultipeSpace(cmd)
	args := strings.Split(final, " ")

	switch strings.ToUpper(args[0]) {
	case GET:
		getAllNamespace("all")
	case CREATE:
		if len(args) < 2 {
			fmt.Println("No Valid NS Name")
		} else {
			usr.Namespace = args[1]
			err := newNS()
			if err != nil {
				fmt.Printf("[Create Namespace Error] [%s]\n", err.Error())
			}
		}
	case REMOVE:
		if len(args) < 2 {
			fmt.Println("No Valid NS Name")
		} else {
			usr.Namespace = args[1]
			err := removeNS()
			if err != nil {
				fmt.Printf("[Remove Namespace Error] [%s]\n", err.Error())
			}
		}
	case QUTI:
		fallthrough
	case "Q":
		fallthrough
	case EXIT:
		ret = true
		return
	default:
		fmt.Println("Invalid Command")
	}
	return
}

func newNS() (err error) {
	if usr.ClusterID == "" {
		err = getClusterID()
		if err != nil{
			return
		}
	}
	path := getAPI("/cloud/namespace/create")
	path += "?clusterid=" + usr.ClusterID + "&name=" + usr.Namespace + "&desc=create-by-EQCloud"
	if debug{
		fmt.Printf("[Create NS] Invoke API [%s]\n",path)
	}
	_, err = SandHttp(http.MethodPost, path)
	if err != nil {
		return
	}

	return
}

func removeNS() (err error) {
	if usr.ClusterID == "" {
		err = getClusterID()
		if err != nil{
			return
		}
	}

	path := getAPI("/cloud/namespace/delete")

	path += "?clusterid=" + usr.ClusterID + "&name=" + usr.Namespace
	if debug{
		fmt.Printf("[Remove NS] Invoke API [%s]\n",path)
	}

	_, err = SandHttp(http.MethodPost, path)
	if err != nil {
		return
	}

	return
}
