package main

import (
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"io/ioutil"
)

func pingAction(c *cli.Context) error {
	_, err := SandHttp(http.MethodGet, getPing(c))
	if err != nil {
		fmt.Printf("Ping [%s] Faild [%s]\n", sip, err.Error())
	} else {
		fmt.Printf("Ping [%s] Succ!\n", sip)
		err := ioutil.WriteFile(getStorePath(), []byte(sip), 0777)
		if err != nil {
			fmt.Println(err)
			//os.Exit(-1)
		}
	}

	return nil
}

func getPing(c *cli.Context) string {
	getSIP(c)
	return "http://" + sip + "/_ping"
}

func getStrictPing(ip string)string{
	return "http://" + ip + "/_ping"
}