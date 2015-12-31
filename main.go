package main

import (
	"flag"
	"fmt"
	"os"
	// "strings"
)

type Arg struct {
	Action       *string
	EtcdAddr     *string
	EtcdNode     *string
	EtcdUser     *string
	EtcdPassword *string
	Recursive    *bool
	Retry        *int64
	BackupFile   *string
	RestoreFile  *string
}

var Args = Arg{}

func usage() {
	fmt.Println("-h to get help message")
	os.Exit(1)
}

func init() {
	Args.Action = flag.String("action", "", "only support backup or restore")
	Args.EtcdAddr = flag.String("etcdaddr", "http://127.0.0.1:2379", "backup etcd address")
	Args.EtcdNode = flag.String("etcdnode", "/", "backup etcd node")
	Args.EtcdUser = flag.String("username", "", "if not auth enable, not need")
	Args.EtcdPassword = flag.String("password", "", "iff not auth enable, not need")
	Args.Recursive = flag.Bool("recursive", true, "recursive backup etcd node")
	Args.BackupFile = flag.String("backupfile", "backup.db", "backup filename")
	Args.RestoreFile = flag.String("restorefile", "", "restore filename")
	flag.Parse()

	switch *Args.Action {
	case "backup":
		if *Args.BackupFile == "" {
			fmt.Println("backup filename can not be null")
			usage()
		}
	case "restore":
		if *Args.RestoreFile == "" {
			fmt.Println("restore filename can not be null")
			usage()
		}
	default:
		usage()
	}

}

func main() {
	if *Args.Action == "backup" {
		BackupEtcd(Args)
	} else {
		// fmt.Printf("are you sure to restore %s, yes/no:", *Args.EtcdNode)
		// var input string
		// fmt.Scanf("%s", &input)
		// if strings.Trim(input, "\n") != "yes" {
		// 	os.Exit(1)
		// }
		RestoreEtcd(Args)
	}
}
