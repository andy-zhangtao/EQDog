package main

import (
	"fmt"
)

// User 客户端用户操作数据
type User struct {
	Sip       string    `json:"sip"`
	Region    string    `json:"region"`
	Namespace string    `json:"namespace"`
	ClusterID string    `json:"cluster_id"`
	Svc       SvcConf   `json:"svc"`
	Con       Container `json:"con"`
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

func containerInfo() {
	fmt.Printf("Name: %s \n", usr.Con.Name)
	fmt.Printf("Service: %s \n", usr.Con.Svc)
	fmt.Printf("Namespace: %s \n", usr.Con.Nsme)
	fmt.Printf("Image: %s \n", usr.Con.Img)
	fmt.Printf("Cmd: %s \n", usr.Con.Cmd)
	fmt.Printf("Env: %s \n", usr.Con.Env)
	fmt.Printf("Idx: %v \n", usr.Con.Idx)
}

func check(kind string) (ret bool) {
	ret = true
	switch kind {
	case "namespace":
		if usr.Namespace == "" {
			fmt.Println("Please Choose Namespace")
			ret = false
			return
		}
	case "region":
		if usr.Region == "" {
			fmt.Println("Please Choose Region")
			ret = false
			return
		}
	case "svc.name":
		if usr.Svc.Name == "" {
			fmt.Println("Please Choose Service")
			ret = false
			return
		}
	case "svc":
		if usr.Svc.Name == "" {
			fmt.Println("Please Type svc name")
			ret = false
			return
		}
		if usr.Svc.Netconf.InPort < 0 || usr.Svc.Netconf.InPort > 65536 {
			fmt.Println("Svc docker port invalid")
			ret = false
			return
		}

		if usr.Svc.Netconf.OutPort < 0 || usr.Svc.Netconf.OutPort > 65536 {
			fmt.Println("Svc lb port invalid")
			ret = false
			return
		}
	case "con.id":
		if usr.Con.ID == "" {
			fmt.Println("Please Choose Container")
			ret = false
			return
		}
	case "con":
		if usr.Con.Svc == "" {
			usr.Con.Svc = usr.Svc.Name
		}
		if usr.Con.Nsme == "" {
			usr.Con.Nsme = usr.Namespace
		}
		if usr.Con.Name == "" {
			fmt.Println("Please Type Container Name")
			ret = false
			return
		}
		if usr.Con.Img == "" {
			fmt.Println("Please Type Container Img")
			ret = false
			return
		}
	}

	return
}
