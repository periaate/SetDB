package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"setdb/data"
	"strings"
)

func main() {
	// Passing arguments to this is interpreted as Junction
	// This can be used a a lambda to other processses

	server := flag.Bool("r", false, "run as server")
	shell := flag.Bool("s", false, "run as shell")
	remoteShell := flag.Bool("rs", false, "run as shell")
	command := flag.Bool("c", false, "run a single command")

	flag.Parse()

	data.Fp = "c:/users/daniel/go/src/setdb/storage"
	d := data.Data{}
	d.Sets.Init(10)
	d.Els.Init(10)
	d.Generate()

	if *server {
		data.Serve(&d)
	}

	if *shell {
		Listen(&d)
	}

	if *remoteShell {
		listenRemote(&d)
	}

	if *command {
		var res []string
		d.Command(os.Args[2:], &res)
		for _, line := range res {
			fmt.Println(line)
		}
	}
}

// Listen listens stdin for commands
func Listen(d *data.Data) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		cmd := scanner.Text()
		args := strings.Split(cmd, " ")
		var res []string
		d.Command(args, &res)
		for _, line := range res {
			fmt.Println(line)
		}
	}
}

func listenRemote(d *data.Data) {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		cmd := scanner.Text()
		args := strings.Split(cmd, " ")
		var res []string
		err = client.Call("Data.Command", args, &res)
		if err != nil {
			log.Fatal("error:", err)
		}

		for _, line := range res {
			fmt.Println(line)
		}
	}
}
