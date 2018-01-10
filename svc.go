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

type HttpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// SvcConf 服务配置信息
// 默认情况下Replicas为1
type SvcConf struct {
	Id        string       `json:"_id,omitempty"`
	Name      string       `json:"name"`
	Desc      string       `json:"desc"`
	Replicas  int          `json:"replicas"`
	Namespace string       `json:"namespace"`
	Netconf   NetConfigure `json:"netconf"`
}

// NetConfigure 服务配置信息
// accessType 默认为ClusterIP:
//     0 - ClusterIP
//     1 - LoadBalancer
//     2 - SvcLBTypeInner
// Inport 容器监听端口
// Outport 负载监听端口
// protocol 协议类型 默认为TCP
//     0 - TCP
//     1 - UDP
type NetConfigure struct {
	AccessType int `json:"access_type"`
	InPort     int `json:"in_port"`
	OutPort    int `json:"out_port"`
	Protocol   int `json:"protocol"`
}

// SvcConfGroup 服务群组配置信息
// 作为自己的软服务编排(以业务场景为主,进行的服务编排.不依赖于k8s的服务编排)
type SvcConfGroup struct {
	SvcGroup  map[string]int `json:"svc_group"`
	Namespace string         `json:"namespace"`
	Clusterid string         `json:"clusterid"`
	Name      string         `json:"name"`
}

func svcOperation() {
	for {
		prompt := promptui.Prompt{
			Label:     "EQCloud Service >",
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

		if ret := switchSVC(result); ret {
			return
		}
	}
}

func switchSVC(cmd string) (ret bool) {

	ret = false

	final := gs.RemoveMultipeSpace(cmd)
	args := strings.Split(final, " ")

	switch strings.ToUpper(args[0]) {
	case GET:
		err := getAllSvc()
		if err != nil {
			fmt.Printf("[Get Service Error] [%s]\n", err.Error())
		}
	case NEW:
		if len(args) < 2 {
			fmt.Println("No Valid Svc Name")
		} else {
			//usr.Namespace = args[1]
			err := newSvc(strings.TrimSpace(args[1]))
			if err != nil {
				fmt.Printf("[Create Service Error] [%s]\n", err.Error())
			}
		}
	case CREATE:
		if len(args) < 2 {
			fmt.Println("No Valid Svc Name")
		} else {
			err := createSvc(strings.TrimSpace(args[1]))
			if err != nil {
				fmt.Printf("[Create Service Error] [%s]\n", err.Error())
			}
		}
	case CHECK:
		check("svc")
	case REMOVE:
		removeSvc()
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

func getAllSvc() (err error) {
	if usr.Namespace == "" {
		fmt.Println("Choose Namespace First")
		return
	}
	path := getAPI("/cloud/svcconf/info")
	path += "?namespace=" + usr.Namespace

	if debug {
		fmt.Printf("[Get ALl SVC] Invoke API [%s] \n", path)
	}

	data, err := SandHttp(http.MethodGet, path)
	if err != nil {
		return
	}

	var vss []SvcConf

	if debug {
		fmt.Printf("[Get ALl SVC] Response Data [%s] \n", string(data))
	}
	err = json.Unmarshal(data, &vss)
	if err != nil {
		return
	}

	var vsp []string

	for _, s := range vss {
		vsp = append(vsp, s.Name)
	}

	cns, err := prompt(SVCLabel, vsp)
	if err != nil {
		fmt.Println(err)
	}

	for _, s := range vss {
		if strings.Compare(cns, s.Name) == 0 {
			usr.Svc = s
		}
	}

	return
}

func newSvc(name string) (err error) {

	svc := SvcConf{}
	svc.Name = name

	//validate := func(input string) error {
	//	//_, err := strconv.Atoi(input)
	//	//if err != nil {
	//	//	fmt.Printf("[Create Svc] Error. This Value [%v] Must Be Number Type \n", input)
	//	//}
	//
	//	return nil
	//}

	prompt := promptui.Prompt{
		Label:     "EQCloud Service Desc >",
		IsVimMode: true,
		AllowEdit: true,
	}

	svc.Desc, err = prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
			return
		}
		fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
		os.Exit(-1)
	}

	for {
		if svc.Replicas != 0 {
			break
		}
		reprompt := promptui.Prompt{
			Label:     "EQCloud Service Replicas >",
			IsVimMode: true,
			AllowEdit: true,
			//Validate:  validate,
			Default: "1",
		}

		replicas, err := reprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				break
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		svc.Replicas, _ = strconv.Atoi(replicas)

	}

	var net NetConfigure
	atprompt := promptui.Prompt{
		Label:     "EQCloud Service AccessType >",
		IsVimMode: true,
		AllowEdit: true,
		//Validate:  validate,
		Default: "0",
	}

	replicas, err := atprompt.Run()
	if err != nil {
		if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
			return
		}
		fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
		os.Exit(-1)
	}

	net.AccessType, _ = strconv.Atoi(replicas)

	//prompt.Validate = nil
	//prompt.Label = "EQCloud Service Protocol >"
	prompt = promptui.Prompt{
		Label:     "EQCloud Service Protocol >",
		IsVimMode: true,
		AllowEdit: true,
		//Validate:  validate,
		Default: "0",
	}

	replicas, err = prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
			return
		}
		fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
		os.Exit(-1)
	}

	net.Protocol, _ = strconv.Atoi(replicas)

	for {
		if net.InPort != 0 {
			break
		}

		dpprompt := promptui.Prompt{
			Label:     "EQCloud Service Docker Port >",
			IsVimMode: true,
			AllowEdit: true,
			//Validate:  validate,
		}

		replicas, err = dpprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				return
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		net.InPort, _ = strconv.Atoi(replicas)
	}

	for {
		if net.OutPort != 0 {
			break
		}

		lbprompt := promptui.Prompt{
			Label:     "EQCloud Service LB Port >",
			IsVimMode: true,
			AllowEdit: true,
			//Validate:  validate,
		}

		replicas, err = lbprompt.Run()
		if err != nil {
			if err == promptui.ErrAbort || err == promptui.ErrEOF || err == promptui.ErrInterrupt {
				return
			}
			fmt.Printf("[EQCloud] Meet Error[%s]\n", err.Error())
			os.Exit(-1)
		}

		net.OutPort, _ = strconv.Atoi(replicas)
	}

	svc.Netconf = net

	usr.Svc = svc
	return
}

func createSvc(name string) (err error) {
	if debug {
		fmt.Printf("[Create Svc] The Svc Name [%s] Current Svc Name [%s] \n", name, usr.Svc.Name)
	}

	if usr.Svc.Namespace == "" {
		usr.Svc.Namespace = usr.Namespace
	}
	if ret := check("namespace"); !ret {
		return
	}

	if ret := check("svc"); !ret {
		return
	}

	path := getAPI("/cloud/svcconf/check")
	ds, _ := json.Marshal(usr.Svc)

	body, err := SandHttpBody(http.MethodPost, path, string(ds))
	if err != nil {
		fmt.Printf("[Create Svc] Error [%s]\n", err.Error())
		return
	}

	var ce HttpError

	err = json.Unmarshal(body, &ce)
	if err != nil {
		fmt.Printf("[Create Svc] Error [%s] [%s]\n", err.Error(), string(body))
		return
	}

	if ce.Code == 1001 {
		fmt.Println("Create Svc Succ")
	} else {
		fmt.Printf("Create Svc Failed [%d] \n", ce.Code)
	}
	return
}

func removeSvc() (err error) {

	if ret := check("namespace"); !ret {
		return
	}

	path := getAPI("/cloud/svcconf/info")
	path += "?namespace=" + usr.Namespace

	if debug {
		fmt.Printf("[Get ALl SVC] Invoke API [%s] \n", path)
	}

	data, err := SandHttp(http.MethodGet, path)
	if err != nil {
		return
	}

	var vss []SvcConf

	if debug {
		fmt.Printf("[Get ALl SVC] Response Data [%s] \n", string(data))
	}
	err = json.Unmarshal(data, &vss)
	if err != nil {
		return
	}

	var vsp []string

	for _, s := range vss {
		vsp = append(vsp, s.Name)
	}

	cns, err := prompt(SVCLabel, vsp)
	if err != nil {
		fmt.Println(err)
	}

	var id string
	var name string
	for _, s := range vss {
		if strings.Compare(cns, s.Name) == 0 {
			id = s.Id
			name = s.Name
		}
	}

	prompt := promptui.Prompt{
		Label:     "Delete Service " + name,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if strings.Compare(strings.ToUpper(strings.TrimSpace(result)), "Y") == 0 {
		path = getAPI("/cloud/svcconf/delete")
		path += "?namespace=" + usr.Namespace + "&id=" + id
		if debug {
			fmt.Printf("[Rm SVC] Invoke API [%s] \n", path)
		}

		_, err = SandHttp(http.MethodPost, path)
		if err != nil {
			fmt.Printf("[Rm SVC] Invoke API Failed [%s] \n", err.Error())
			return
		}
	}

	return

}
