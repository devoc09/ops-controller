package instance

import (
	"encoding/json"
	"fmt"
	"os"
)

type Instance struct {
	Instance string   `json:"instance"`
	Image    string   `json:"image"`
	Ports    []string `json:"ports"`
}

type ControllerInstance struct {
	Instance string   `json:"instance"`
	Image    string   `json:"image"`
	Ports    []string `json:"ports"`
	AutoHeal bool     `json:"autoheal"`
}

func New(abs string) (i *Instance) {
	f, _ := os.Open(abs)
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&i); err != nil {
		fmt.Printf("Decode json file Error\n")
		panic(err)
	}
	return
}
