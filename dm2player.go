package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

type Config struct {
	Baseq2 string `json:"baseq2folder"`
	Q2exe  string `json:"q2binary"`
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	sep := string(os.PathSeparator)
	configfile := fmt.Sprintf("%s%sdm2player.json", user.HomeDir, sep)

	configbody, err := os.ReadFile(configfile)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(configbody, &config)
	if err != nil {
		panic(err)
	}

	// copy the demo the right place
	demoname := fmt.Sprintf("%s%s%s%stempdemo.dm2", config.Baseq2, sep, "demos", sep)
	demosrc, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(demoname, demosrc, 0666)
	if err != nil {
		panic(err)
	}

	// make a temporary config and write it to baseq2 folder
	cfg := "alias loopdemo \"disconnect; demo tempdemo; set nextserver loopdemo\""
	cfgname := fmt.Sprintf("%s%stempdemo.cfg", config.Baseq2, sep)
	err = os.WriteFile(cfgname, []byte(cfg), 0666)
	if err != nil {
		panic(err)
	}

	//cmd := fmt.Sprintf("%s +exec tempdemo.cfg", config.Q2exe)
	cmd := exec.Command(config.Q2exe, "+demo tempdemo")
	_ = cmd.Run()

	err = os.Remove(demoname)
	if err != nil {
		panic(err)
	}

	err = os.Remove(cfgname)
	if err != nil {
		panic(err)
	}
}
