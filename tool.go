package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"strings"
	"os/user"
	"io/ioutil"
	"github.com/manifoldco/promptui"
)

const (
	NSLabel  = "Choose The Namespace"
	RGLabel  = "Choose The Region"
	CLULabel = "Choose The Cluster ID"
	SVCLabel = "Choose The Svc"
	CONLabel = "Choose The Container"
)

func getAPI(path string) string {
	return "http://" + sip + "/v1" + path
}

func getSIP(c *cli.Context) {
	if !sipParse {
		sip = c.GlobalString(SIP)
		oip := getStoreEndpoint()

		if !strings.Contains(sip, oip) && strings.Compare(sip, "127.0.0.1:8000") == 0 {
			sip = oip
			sipParse = true
		}
	}
}

// getStoreEndpoint 获取已经存储的服务端IP
func getStoreEndpoint() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	path := ""
	if !strings.HasSuffix(usr.HomeDir, "/") {
		path = usr.HomeDir + "/.tdog"
	} else {
		path = usr.HomeDir + ".tdog"
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ""
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return string(data)
}

// getStorePath 获取配置文件路径
func getStorePath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	path := ""
	if !strings.HasSuffix(usr.HomeDir, "/") {
		path = usr.HomeDir + "/.tdog"
	} else {
		path = usr.HomeDir + ".tdog"
	}

	return path
}

func prompt(label string, items []string) (result string, err error) {
	prompt := promptui.Select{
		//Label: "Select Day",
		//Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
		//	"Saturday", "Sunday"},
		Label: label,
		Items: items,
		Size:  10,
	}

	_, result, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	return
}
