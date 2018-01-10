package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/manifoldco/promptui"
	"strings"
)

const (
	SIP = "server"
)

var sip string
var sipParse bool
var usr *User
var debug bool

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help"}
	usr = new(User)
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
		cli.BoolFlag{
			Name:  "debug",
			Usage: "More Debug Info",
		},
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
			Name:        "region",
			Aliases:     []string{"rg"},
			Usage:       "EQCloud Region",
			Description: "Choose The Region You Want To Deploy",
			Action:      metaAction,
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

	app.Run(os.Args)
}

func cliAction(c *cli.Context) error {
	debug = c.GlobalBool("debug")
	getSIP(c)

	err := ioutil.WriteFile(getStorePath(), []byte(sip), 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for {
		prompt := promptui.Prompt{
			Label:     "EQCloud >",
			IsVimMode: true,
			AllowEdit: true,
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				os.Exit(0)
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}
		switch strings.ToUpper(result) {
		case PING:
			pingAction(c)
		case REGION:
			metaAction(c)
		case INFO:
			configure()
		case NAMESPACE:
			nsOperation()
		case SVC:
			svcOperation()
		case CONTAINER:
			containerOperation()
		case QUTI:
			fallthrough
		case "Q":
			fallthrough
		case EXIT:
			fmt.Println("Bye!")
			os.Exit(0)
		}
	}

	return nil
}
