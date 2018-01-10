package main

import (
	"github.com/manifoldco/promptui"
	"fmt"
	"os"
	gs "github.com/andy-zhangtao/gogather/strings"

	"strings"
	"net/http"
	"encoding/json"
	"strconv"
)

type Container struct {
	ID   string            `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Cmd  []string          `json:"cmd"`
	Env  map[string]string `json:"env"`
	Svc  string            `json:"svc"`
	Nsme string            `json:"namespace"`
	Idx  int               `json:"idx"`
}

func containerOperation() {
	for {
		prompt := promptui.Prompt{
			Label:     "EQCloud Container >",
			IsVimMode: true,
			AllowEdit: true,
		}

		cmd, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		if ret := switchCon(cmd); ret {
			return
		}
	}
}

func switchCon(cmd string) (ret bool) {

	ret = false

	final := gs.RemoveMultipeSpace(cmd)
	args := strings.Split(final, " ")

	switch strings.ToUpper(args[0]) {
	case GET:
		err := getAllCon()
		if err != nil {
			fmt.Printf("[Get Service Error] [%s]\n", err.Error())
		}
	case CAT:
		containerInfo()
	case NEW:
		if len(args) < 2 {
			fmt.Println("No Valid Container Name")
		} else {
			err := newCon(strings.TrimSpace(args[1]))
			if err != nil {
				fmt.Printf("[Create Container Error] [%s]\n", err.Error())
			}
		}
	case CREATE:
		if len(args) < 2 {
			fmt.Println("No Valid Container Name")
		} else {
			err := createCon(strings.TrimSpace(args[1]))
			if err != nil {
				fmt.Printf("[Create Container Error] [%s]\n", err.Error())
			}
		}
	case CHECK:
		check("con")
	case REMOVE:
		rmCon()
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

func getAllCon() (err error) {
	if ret := check("namespace"); !ret {
		return
	}

	if ret := check("svc.name"); !ret {
		return
	}

	path := getAPI("/cloud/container/info")
	path += "?namespace=" + usr.Namespace + "&svc=" + usr.Svc.Name

	if debug {
		fmt.Printf("[Get ALl CONTAINER] Invoke API [%s] \n", path)
	}

	data, err := SandHttp(http.MethodGet, path)
	if err != nil {
		return
	}

	var vss []Container

	if debug {
		fmt.Printf("[Get ALl CONTAINER] Response Data [%s] \n", string(data))
	}
	err = json.Unmarshal(data, &vss)
	if err != nil {
		return
	}

	var vsp []string

	for _, s := range vss {
		vsp = append(vsp, s.Name)
	}

	if len(vsp) == 0 {
		fmt.Println("There are no containers")
		return
	}
	cns, err := prompt(CONLabel, vsp)
	if err != nil {
		fmt.Println(err)
	}

	for _, s := range vss {
		if strings.Compare(cns, s.Name) == 0 {
			usr.Con = s
		}
	}

	return
}

func rmCon() (err error) {
	if ret := check("namespace"); !ret {
		return
	}

	if ret := check("svc.name"); !ret {
		return
	}

	path := getAPI("/cloud/container/info")
	path += "?namespace=" + usr.Namespace + "&svc=" + usr.Svc.Name

	if debug {
		fmt.Printf("[Get ALl CONTAINER] Invoke API [%s] \n", path)
	}

	data, err := SandHttp(http.MethodGet, path)
	if err != nil {
		return
	}

	var vss []Container

	if debug {
		fmt.Printf("[Get ALl CONTAINER] Response Data [%s] \n", string(data))
	}
	err = json.Unmarshal(data, &vss)
	if err != nil {
		return
	}

	var vsp []string

	for _, s := range vss {
		vsp = append(vsp, s.Name)
	}

	cns, err := prompt(CONLabel, vsp)
	if err != nil {
		fmt.Println(err)
	}

	var id string
	var name string
	for _, s := range vss {
		if strings.Compare(cns, s.Name) == 0 {
			id = s.ID
			name = s.Name
		}
	}

	prompt := promptui.Prompt{
		Label:     "Delete Container " + name,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if strings.Compare(strings.ToUpper(strings.TrimSpace(result)), "Y") == 0 {
		path = getAPI("/cloud/container/delete")
		path += "?namespace=" + usr.Namespace + "&id=" + id + "&svc=" + usr.Svc.Name
		if debug {
			fmt.Printf("[Rm Con] Invoke API [%s] \n", path)
		}

		_, err = SandHttp(http.MethodPost, path)
		if err != nil {
			fmt.Printf("[Rm Con] Invoke API Failed [%s] \n", err.Error())
			return
		}
	}

	return
}

func newCon(name string) (err error) {
	con := Container{}
	con.Name = name

	if usr.Con.Nsme == "" {
		usr.Con.Nsme = usr.Namespace
	}

	if usr.Con.Svc == "" {
		usr.Con.Svc = usr.Svc.Name
	}

	prompt := promptui.Prompt{
		Label:     "EQCloud Container Img >",
		IsVimMode: true,
		AllowEdit: true,
	}

	con.Img, err = prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
			return
		}
		fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
		os.Exit(-1)
	}

	for {
		atprompt := promptui.Prompt{
			Label:     "EQCloud Container Cmd >",
			IsVimMode: true,
			AllowEdit: true,
			Default:   "",
		}

		cmd, err := atprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		if strings.ToUpper(cmd) == "Q" {
			break
		}

		con.Cmd = append(con.Cmd, cmd)
	}

	ev := make(map[string]string)
	for {
		atprompt := promptui.Prompt{
			Label:     "EQCloud Container Env >",
			IsVimMode: true,
			AllowEdit: true,
			Default:   "",
		}

		env, err := atprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		if strings.ToUpper(env) == "Q" {
			break
		}

		e := strings.Split(env, "=")
		if len(e) == 1 {
			fmt.Printf("Not Support Singal Side ENV")
		} else {
			ev[e[0]] = e[1]
		}
	}

	con.Env = ev
	for {
		if con.Idx != 0 {
			break
		}
		reprompt := promptui.Prompt{
			Label:     "EQCloud Container Idx >",
			IsVimMode: true,
			AllowEdit: true,
			Default:   "1",
		}

		replicas, err := reprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		con.Idx, _ = strconv.Atoi(replicas)

	}

	usr.Con = con

	return
}

func createCon(name string) (err error) {
	if debug {
		fmt.Printf("[Create Container] The Container Name [%s] Current Container Name [%s] \n", name, usr.Svc.Name)
	}

	if usr.Con.Nsme == "" {
		usr.Con.Nsme = usr.Namespace
	}

	if usr.Con.Svc == "" {
		usr.Con.Svc = usr.Svc.Name
	}

	if ret := check("namespace"); !ret {
		return
	}

	if ret := check("svc.name"); !ret {
		return
	}

	path := getAPI("/cloud/container/create")
	ds, _ := json.Marshal(usr.Con)

	_, err = SandHttpBody(http.MethodPost, path, string(ds))
	if err != nil {
		fmt.Printf("[Create Container] Error [%s]\n", err.Error())
		return
	}

	return
}
