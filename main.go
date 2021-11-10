package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/devoc09/ops-wrap/controller"
)

func main() {
	opshome := GetOpsHome()
	instancePath := path.Join(opshome, "instances")

	files, err := ioutil.ReadDir(instancePath)
	if err != nil {
		return
	}

	for _, f := range files {
		fullpath := path.Join(instancePath, f.Name())
		pid, err := strconv.ParseInt(f.Name(), 10, 32)
		if err != nil {
			panic(err)
		}
		process, err := os.FindProcess(int(pid))
		if err != nil {
			panic(err)
		}
		if err = process.Signal(syscall.Signal(0)); err != nil {
			errMsg := strings.ToLower(err.Error())
			if strings.Contains(errMsg, "already finished") ||
				strings.Contains(errMsg, "already released") ||
				strings.Contains(errMsg, "not initialized") {
				os.Remove(fullpath)
				continue
			}
		}
		body, err := ioutil.ReadFile(fullpath)
		if err != nil {
			panic(err)
		}

		var i controller.Instance
		if err := json.Unmarshal(body, &i); err != nil {
			panic(err)
		}
		fmt.Printf("Process ID: %s\n", f.Name())
		fmt.Println(i)
	}

	// Watch Directory
	controller.Watch(instancePath)
}

var homeDir = ""

func HomeDir() (string, error) {
	if homeDir != "" {
		return homeDir, nil
	}

	var err error = nil

	if homeDir, err = os.UserHomeDir(); err == nil {
		return homeDir, nil
	}

	homeDir, err = os.Getwd()
	if err != nil {
		return "", errors.New("home Directory not detected")
	}

	return homeDir, nil
}

func GetOpsHome() string {
	home, err := HomeDir()
	if err != nil {
		panic(err)
	}

	opshome := path.Join(home, ".ops")
	return opshome
}
