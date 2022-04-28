package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"setdb/data"
	"setdb/server"
)

func main() {
	// Passing arguments to this is interpreted as Junction
	// This can be used a a lambda to other processses

	remote := flag.Bool("r", false, "run as server")
	shell := flag.Bool("s", false, "run as shell")
	remoteShell := flag.Bool("rs", false, "run as shell")
	command := flag.Bool("c", false, "run a single command")
	configure := flag.Bool("configure", false, "make configuration file")

	flag.Parse()

	if *configure {
		c := new(cfg)
		c.Save()
		return
	}

	c := new(cfg)
	c.Load()

	data.Fp = c.Fp
	d := data.Data{}
	d.Sets.Init(10)
	d.Els.Init(10)
	d.Generate()

	if *remote {
		server.Serve(&d)
	}

	if *shell {
		server.Listen(&d)
	}

	if *remoteShell {
		server.ListenRemote(&d)
	}

	if *command {
		var res []string
		d.Command(os.Args[2:], &res)
		for _, line := range res {
			fmt.Println(line)
		}
	}
}

type cfg struct {
	Fp string `json:"filePath"`
}

// confing loads json as config
func (c *cfg) Load() {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *cfg) Save() {
	fmt.Print("filepath: ")
	var filepath string
	fmt.Scanln(&filepath)
	c.Fp = filepath
	file, _ := json.MarshalIndent(c, "", " ")
	_ = ioutil.WriteFile("./config.json", file, 0644)
}
