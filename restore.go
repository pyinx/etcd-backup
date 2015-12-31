package main

import (
	// 	"flag"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"io/ioutil"
	"os"
	"strings"
)

func RestoreEtcd(args Arg) {
	if strings.HasSuffix(*args.EtcdAddr, "/") {
		*args.EtcdAddr = strings.TrimRight(*args.EtcdAddr, "/")
	}
	client := etcd.NewClient([]string{*args.EtcdAddr})
	if *args.EtcdUser != "" && *args.EtcdPassword != "" {
		client.SetCredentials(*args.EtcdUser, *args.EtcdPassword)
	}
	restorefile, err := os.Open(*args.RestoreFile)
	if err != nil {
		fmt.Printf("open restorefile %s err: %s\n", *args.RestoreFile, err)
		os.Exit(1)
	}
	defer restorefile.Close()
	data, err := ioutil.ReadAll(restorefile)
	if err != nil {
		fmt.Printf("read restorefile %s err: %s\n", *args.RestoreFile, err)
		os.Exit(1)
	}
	var jsonmap []map[string]interface{}
	err = json.Unmarshal(data, &jsonmap)
	if err != nil {
		fmt.Println("unmarshal data err: %s\n", err)
		os.Exit(1)
	}
	for _, node := range jsonmap {
		// fmt.Printf("%s %s %d %t\n", node["Key"], node["Value"], uint64(node["TTL"].(float64)), node["Dir"].(bool))
		if node["Dir"].(bool) {
			client.Set(node["Key"].(string), node["Value"].(string), uint64(node["TTL"].(float64)))
		}
	}
}
