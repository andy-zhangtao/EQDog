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

func svcorOperation() {
	for {
		prompt := promptui.Prompt{
			Label:     "EQCloud SvcOR >",
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

		if ret := switchSOR(cmd); ret {
			return
		}
	}
}

func switchSOR(cmd string) (ret bool) {

	ret = false

	final := gs.RemoveMultipeSpace(cmd)
	args := strings.Split(final, " ")

	switch strings.ToUpper(args[0]) {
	case GET:
		err := getALlSor()
		if err != nil {
			fmt.Printf("[Get Service Error] [%s]\n", err.Error())
		}
	case CAT:
		svcorInfo()
	case NEW:
		if len(args) < 2 {
			fmt.Println("No Valid SvcOR Name")
		} else {
			err := newSor(strings.TrimSpace(args[1]))
			if err != nil {
				fmt.Printf("[Create SvcOR Error] [%s]\n", err.Error())
			}
		}
	case CREATE:
		if len(args) < 2 {
			fmt.Println("No Valid SvcOR Name")
		} else {
			err := createSor(strings.TrimSpace(args[1]))
			if err != nil {
				fmt.Printf("[Create SvcOR Error] [%s]\n", err.Error())
			}
		}
	case CHECK:
		check("con")
	case REMOVE:
		rmSor()
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

func getALlSor() (err error) {
	if ret := check("namespace"); !ret {
		return
	}

	path := getAPI("/cloud/svc/group/info")
	path += "?namespace=" + usr.Namespace

	if debug {
		fmt.Printf("[Get ALl SVCOR INFO] Invoke API [%s] \n", path)
	}

	data, err := SandHttp(http.MethodGet, path)
	if err != nil {
		return
	}

	var vss HttpResp

	if debug {
		fmt.Printf("[Get ALl SVCOR INFO] Response Data [%s] \n", string(data))
	}

	err = json.Unmarshal(data, &vss)
	if err != nil {
		return
	}

	var vs []SvcConfGroup

	str, _ := vss.Data.(string)
	err = json.Unmarshal([]byte(str), &vs)
	if err != nil {
		return
	}

	var vsp []string

	for _, s := range vs {
		vsp = append(vsp, s.Name)
	}

	if len(vsp) == 0 {
		fmt.Println("There are no SVCOR")
		return
	}
	cns, err := prompt(SVCLabel, vsp)
	if err != nil {
		fmt.Println(err)
	}

	for _, s := range vs {
		if strings.Compare(cns, s.Name) == 0 {
			usr.SvcOR = s
		}
	}

	return
}
func newSor(name string) (err error) {
	svc := SvcConfGroup{}
	svc.Name = name

	svc.Clusterid = usr.ClusterID
	svc.Namespace = usr.Namespace

	sg := make(map[string]int)

	for {

		reprompt := promptui.Prompt{
			Label:     "EQCloud SvcOR SvcName >",
			IsVimMode: true,
			AllowEdit: true,
			//Validate:  validate,
			//Default: "1",
		}

		name, err := reprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		if name == "" {
			fmt.Println("SvcName Cannot Be Empty")
			continue
		}

		if strings.ToUpper(name) == "Q" {
			break
		}

		reprompt = promptui.Prompt{
			Label:     "EQCloud SvcOR Idx >",
			IsVimMode: true,
			AllowEdit: true,
			//Validate:  validate,
			Default: "1",
		}

		idx, err := reprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		dx, err := strconv.Atoi(idx)
		if err != nil {
			dx = 1
		}
		sg[name] = dx
	}

	svc.SvcGroup = sg

	usr.SvcOR = svc
	return
}

func createSor(name string) (err error) {
	if debug {
		fmt.Printf("[Create SvcOR] The SvcOR Name [%s] Current SvcOR Name [%s] \n", name, usr.SvcOR.Name)
	}

	if usr.SvcOR.Namespace == "" {
		usr.SvcOR.Namespace = usr.Namespace
	}

	if ret := check("namespace"); !ret {
		return
	}

	if usr.ClusterID == "" {
		err = getClusterID()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println(usr)
	path := getAPI("/cluod/svc/group/add")
	path += "?name=" + usr.SvcOR.Name + "&clusterid=" + usr.ClusterID + "&namespace=" + usr.Namespace
	oldPath := path
	for key, value := range usr.SvcOR.SvcGroup {
		path = oldPath
		path += "&svcname=" + key + "&idx=" + strconv.Itoa(value)
		_, err := SandHttp(http.MethodPost, path)
		if err != nil {
			fmt.Printf("[Create SvcOR] Error [%s]\n", err.Error())
			break
		}
	}

	return
}

func rmSor() {
	if ret := check("namespace"); !ret {
		return
	}

	path := getAPI("/cloud/svc/group/info")
	path += "?namespace=" + usr.Namespace

	if debug {
		fmt.Printf("[Get ALl SVCOR INFO] Invoke API [%s] \n", path)
	}

	data, err := SandHttp(http.MethodGet, path)
	if err != nil {
		return
	}

	var vss HttpResp

	if debug {
		fmt.Printf("[Get ALl SVCOR INFO] Response Data [%s] \n", string(data))
	}

	err = json.Unmarshal(data, &vss)
	if err != nil {
		return
	}

	var vs []SvcConfGroup

	str, _ := vss.Data.(string)
	err = json.Unmarshal([]byte(str), &vs)
	if err != nil {
		return
	}

	var vsp []string

	for _, s := range vs {
		vsp = append(vsp, s.Name)
	}

	if len(vsp) == 0 {
		fmt.Println("There are no SvcOR")
		return
	}
	cns, err := prompt(SVCLabel, vsp)
	if err != nil {
		fmt.Println(err)
	}

	//var id string
	var name string
	for _, s := range vs {
		if strings.Compare(cns, s.Name) == 0 {
			//id = s.ID
			name = s.Name
		}
	}

	prompt := promptui.Prompt{
		Label:     "Delete SvcOR " + name,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Cancel")
		return
	}

	if strings.Compare(strings.ToUpper(strings.TrimSpace(result)), "Y") == 0 {
		path = getAPI("/cloud/svc/group/delete")
		path += "?namespace=" + usr.Namespace + "&name=" + usr.SvcOR.Name
		if debug {
			fmt.Printf("[Rm SvcOR] Invoke API [%s] \n", path)
		}

		_, err = SandHttp(http.MethodPost, path)
		if err != nil {
			fmt.Printf("[Rm SvcOR] Invoke API Failed [%s] \n", err.Error())
			return
		}
	}

	return
}
