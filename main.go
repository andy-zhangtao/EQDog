package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	SIP = "server"
)

var sip string

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help"}
}

func main() {
	app := cli.NewApp()

	app.Name = "eqdog"
	app.Usage = "Make Deploy More Easy!"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "ZhangTao",
			Email: "ztao@gmail.com",
		},
	}
	app.Copyright = "(c) 2018 By Andy Zhang"
	app.EnableBashCompletion = true

	app.Action = cliAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  SIP,
			Value: "127.0.0.1:8000",
			Usage: "The DDog Server IP",
		},
		//cli.BoolFlag{
		//	Name:  "ping",
		//	Usage: "Ping DDog Server",
		//},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:            "ping",
			Aliases:         []string{"pg"},
			Usage:           "ping the ddog server",
			Description:     "Check the connect whether success",
			SkipFlagParsing: true,
			Action:          pingAction,
		},
		cli.Command{
			Name:        "namespace",
			Aliases:     []string{"ns"},
			Usage:       "Get/Modify Namespace",
			Description: "You can get/create/modify/delete namespace",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "get",
					Usage: "Get Namespace Info",
				},
			},
			Action: nsAction,
		},
	}

	//app.CommandNotFound = func(c *cli.Context, command string) {
	//	fmt.Printf("Thar be no %q here.\n", command)
	//}
	//app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
	//	if isSubcommand {
	//		return err
	//	}
	//
	//	fmt.Printf("WRONG: %#v\n", err)
	//	return nil
	//}

	app.Run(os.Args)
}

func cliAction(c *cli.Context) error {

	getSIP(c)

	err := ioutil.WriteFile(getStorePath(), []byte(sip), 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//sip = c.GlobalString(SIP)
	//if sip == "" {
	//	if _, err := os.Stat(path); os.IsNotExist(err) {
	//		fmt.Println("Please Specify Server IP!")
	//		os.Exit(-1)
	//	} else {
	//		data, err := ioutil.ReadFile(path)
	//		if err != nil {
	//			fmt.Println(err)
	//			os.Exit(-1)
	//		}
	//		sip = string(data)
	//	}
	//} else {
	//	err := ioutil.WriteFile(path, []byte(sip), 0777)
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(-1)
	//	}
	//}

	return nil
}

func pingAction(c *cli.Context) error {
	_, err := SandHttp(http.MethodGet, getPing(c))
	if err != nil {
		fmt.Printf("Ping [%s] Faild [%s]\n", sip, err.Error())
	} else {
		fmt.Printf("Ping [%s] Succ!", sip)
		err := ioutil.WriteFile(getStorePath(), []byte(sip), 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	return nil
}

//func getAPI(path string) string {
//	getSIP()
//	return "http://" + sip + "/v1" + path
//}

func getPing(c *cli.Context) string {
	getSIP(c)
	return "http://" + sip + "/_ping"
}

func getSIP(c *cli.Context) {
	sip = c.GlobalString(SIP)
	oip := getStoreEndpoint()

	if !strings.Contains(sip, oip) && strings.Compare(sip, "127.0.0.1:8000") == 0 {
		sip = oip
	}
}

func nsAction(c *cli.Context) error {

	return nil
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
