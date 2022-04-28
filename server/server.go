package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"setdb/data"
	"strings"
)

// Serve rpc at port
func Serve(d *data.Data) {
	rpc.Register(d)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "localhost:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
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

// ListenRemote listens and forwards to remote shell
func ListenRemote(d *data.Data) {
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
