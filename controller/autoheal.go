package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devoc09/ops-wrap/controller/instance"
)

func healInstance(src string) error {
	failInstance := getFileName(src)
	ctrpath := controllerHome()
	configpath := filepath.Join(ctrpath, "instances", failInstance)
	config, err := os.Open(configpath)
	if err != nil {
		return fmt.Errorf("Error Open file .ops-controller/instances/ : %w", err)
	}
	defer config.Close()

	decoder := json.NewDecoder(config)
	var ctri instance.ControllerInstance
	if err := decoder.Decode(&ctri); err != nil {
		return fmt.Errorf("Error decode ControllerInstance config file: %w", err)
	}
	if ctri.AutoHeal {
		err := exec.Command("ops", "instance", "create", "ops-hello", "-p", "8080", "ops-hello.img").Start()
		if err != nil {
			return fmt.Errorf("Exec Ops instance create Error: %w", err)
		}
	}
	fmt.Println("Instance Auto Heal Succeeded")
	return nil
}
