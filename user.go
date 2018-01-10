package main

import (
	"fmt"
)

// User 客户端用户操作数据
type User struct {
	Sip       string  `json:"sip"`
	Region    string  `json:"region"`
	Namespace string  `json:"namespace"`
	ClusterID string  `json:"cluster_id"`
	Svc       SvcConf `json:"svc"`
}

func configure() {
	fmt.Printf("Server: %s \n", usr.Sip)
	fmt.Printf("Region: %s \n", usr.Region)
	fmt.Printf("ClusterID: %s \n", usr.ClusterID)
	fmt.Printf("Namespace: %s \n", usr.Namespace)
	fmt.Println("Service: ")
	fmt.Printf("  |----Name: %s\n", usr.Svc.Name)
	fmt.Printf("  |----Replicas: %v\n", usr.Svc.Replicas)
	fmt.Printf("  |----Desc: %s\n", usr.Svc.Desc)
	fmt.Printf("  |------AccessType: %v\n", usr.Svc.Netconf.AccessType)
	fmt.Printf("  |------OutPort: %v\n", usr.Svc.Netconf.OutPort)
	fmt.Printf("  |------InPort: %v\n", usr.Svc.Netconf.InPort)
	fmt.Printf("  |------Protocol: %v\n", usr.Svc.Netconf.Protocol)
}

func check(kind string)(ret bool) {
	ret = true
	switch kind{
	case "namespace":
		if usr.Namespace == ""{
			fmt.Println("Please Choose Namespace")
			ret = false
			return
		}
	case "region":
		if usr.Region == ""{
			fmt.Println("Please Choose Region")
			ret = false
			return
		}
	case "svc":
		if usr.Svc.Name == ""{
			fmt.Println("Please Type svc name")
			ret = false
			return
		}
		if usr.Svc.Netconf.InPort <0 || usr.Svc.Netconf.InPort>65536 {
			fmt.Println("Svc docker port invalid")
			ret = false
			return
		}

		if usr.Svc.Netconf.OutPort <0 || usr.Svc.Netconf.OutPort>65536 {
			fmt.Println("Svc lb port invalid")
			ret = false
			return
		}
	}

	return
}