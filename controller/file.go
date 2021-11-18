package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devoc09/ops-wrap/controller/instance"
)

func homeDir(home string) (string, error) {
	if home != "" {
		return home, nil
	}

	var err error = nil

	if home, err = os.UserHomeDir(); err == nil {
		return home, nil
	}

	home, err = os.Getwd()
	if err != nil {
		return "", errors.New("home Directory not detected")
	}
	return home, nil
}

func controllerHome() string {
	home, err := homeDir("")
	if err != nil {
		panic(err)
	}

	ctrhome := filepath.Join(home, ".ops-controller")
	return ctrhome
}

func createCtrInstanceDir() error {
	ctrhome := controllerHome()

	if err := os.MkdirAll(filepath.Join(ctrhome, "instances"), 0777); err != nil {
		return fmt.Errorf("make ControllerDir Error: %w", err)
	}
	return nil
}

func existDir(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return true
}

func getFileName(path string) string {
	return filepath.Base(path)
}

func writeCtrInstanceFile(src string) error {
	dstname := getFileName(src)
	srcfile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Open File Error: %w\n", err)
	}
	defer srcfile.Close()

	ctrpath := controllerHome()
	dstpath := filepath.Join(ctrpath, "instances")
	if !existDir(dstpath) {
		err := createCtrInstanceDir()
		if err != nil {
			return fmt.Errorf("copyFile createCtrInstanceDir() Error: %w", err)
		}
	}
	dst, err := os.Create(filepath.Join(dstpath, dstname))
	if err != nil {
		return fmt.Errorf("Error Create dst file to .ops-controller/instances/ : %w\n", err)
	}
	defer dst.Close()

	// write controllerInstance info to dst-file.
	decoder := json.NewDecoder(srcfile)
	// var i Instance
	var i instance.Instance
	if err := decoder.Decode(&i); err != nil {
		return fmt.Errorf("Instance json file Decode Error: %w", err)
	}
	encoder := json.NewEncoder(dst)
	// ctri := ControllerInstance{i.Instance, i.Image, i.Ports, true}
	ctri := instance.ControllerInstance{i.Instance, i.Image, i.Ports, true}
	if err := encoder.Encode(ctri); err != nil {
		return fmt.Errorf("Error ControllerInstance json Encode to .ops-controller/instances/ : %w", err)
	}
	// _, err = io.Copy(dst, srcfile)
	// if err != nil {
	// 	return fmt.Errorf("io.Copy() Errof: %w\n", err)
	// }
	return nil
}

// func getPID(abs string) string {
// 	return getFileName(abs)
// }

// func getTargetProcessAbs(abs string) (targetAbs string) {
// 	targetAbs = filepath.Join("/", "proc", getPID(abs))
// 	return
// }

// type Instance struct {
// 	Instance string   `json:"instance"`
// 	Image    string   `json:"image"`
// 	Ports    []string `json:"ports"`
// }

// type ControllerInstance struct {
// 	Instance string   `json:"instance"`
// 	Image    string   `json:"image"`
// 	Ports    []string `json:"ports"`
// 	AutoHeal bool     `json:"autoheal"`
// }
