package main

import (
	// 	"flag"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"strings"
)

type BackupData struct {
	Key   string
	Value string
	TTL   int64
	Dir   bool
}

var backupdata = []BackupData{}

func BackupEtcd(args Arg) {
	if strings.HasSuffix(*args.EtcdAddr, "/") {
		*args.EtcdAddr = strings.TrimRight(*args.EtcdAddr, "/")
	}
	client := etcd.NewClient([]string{*args.EtcdAddr})
	if *args.EtcdUser != "" && *args.EtcdPassword != "" {
		client.SetCredentials(*args.EtcdUser, *args.EtcdPassword)
	}
	data, err := client.Get(*args.EtcdNode, false, *args.Recursive)
	defer client.Close()
	if err != nil {
		fmt.Printf("connect etcd %s err: %s\n", *args.EtcdAddr, err)
		os.Exit(1)
	}
	getNode(data.Node.Nodes)
	backupfile, err := os.OpenFile(*args.BackupFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer backupfile.Close()
	if err != nil {
		fmt.Printf("open backupfile %s err: %s\n", *args.BackupFile, err)
		os.Exit(1)
	}
	out, err := json.Marshal(backupdata)
	if err != nil {
		fmt.Println("marshal data err: %s\n", err)
		os.Exit(1)
	}
	backupfile.Write(out)
}

func getNode(nodes etcd.Nodes) {
	for _, node := range nodes {
		fmt.Println(node.Key)
		if node.Dir {
			backupdata = append(backupdata, BackupData{
				Key:   node.Key,
				Value: node.Value,
				TTL:   node.TTL,
				Dir:   true,
			})
			getNode(node.Nodes)
		} else {
			backupdata = append(backupdata, BackupData{
				Key:   node.Key,
				Value: node.Value,
				TTL:   node.TTL,
				Dir:   false,
			})
		}
	}
}
