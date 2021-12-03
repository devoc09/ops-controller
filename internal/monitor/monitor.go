package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devoc09/ops-wrap/controller/instance"
	"github.com/devoc09/ops-wrap/internal/file"
)

type Monitor struct {
	targetDir string
	instance  *instance.Instance
	Msg       chan instance.Instance // to communicate with the controller
	AutoHeal  bool
}

func New(abs string) (m *Monitor) {
	m = &Monitor{
		targetDir: abs,
		instance:  instance.New(abs),
		Msg:       make(chan instance.Instance),
		AutoHeal:  true,
	}
	return
}

//process to be monitored every second
func (m *Monitor) Start(abs string) {
	go func() {
		for {
			// fmt.Printf("Timer 2 second per Start........\n")
			ok := aliveProcess(filepath.Base(m.targetDir))
			if !ok {
				fmt.Printf("PID: %s Dead!!!\n", filepath.Base(m.targetDir))
				fmt.Printf("Send Channel Dead Instance info\n")
				fmt.Println(*m.instance)
				m.Msg <- *m.instance
				break
				// } else {
				// 	fmt.Printf("PID: %s Alive\n", filepath.Base(abs))
				// }
			}
		}
	}()

	select {
	case hi := <-m.Msg:
		fmt.Printf("channel received......\n")
		fmt.Println(hi)
		fmt.Printf("Monitor Instance Info.........\n")
		fmt.Println(m)
		img := filepath.Base(hi.Image)
		if isAutoHeal(filepath.Base(m.targetDir)) {
			if len(hi.Ports) == 0 {
				if err := exec.Command("ops", "instance", "create", img).Start(); err != nil {
					fmt.Println(err)
				}
			} else {
				if err := exec.Command("ops", "instance", "create", img, "-p", hi.Ports[0]).Start(); err != nil {
					fmt.Println(err)
				}
			}
			fmt.Println("instance heal Succeeded!!")
			return
		}
		return
	}
}

func (m *Monitor) CreateMonitorFile(abs string) error {
	fname := filepath.Base(abs)
	home, err := file.HomeDir("")
	if err != nil {
		panic(err)
	}
	ctrhome := filepath.Join(home, ".ops-controller")
	dstDir := filepath.Join(ctrhome, "instances")
	if !file.ExitDir(dstDir) {
		if err := os.MkdirAll(dstDir, 0777); err != nil {
			panic(err)
		}
	}
	dst, err := os.Create(filepath.Join(dstDir, fname))
	if err != nil {
		return fmt.Errorf("Error Create Monitor File to .ops-controller/instances/ %w\n", err)
	}
	defer dst.Close()
	encoder := json.NewEncoder(dst)
	ctri := instance.ControllerInstance{
		Instance: m.instance.Instance,
		Image:    m.instance.Image,
		Ports:    m.instance.Ports,
		AutoHeal: m.AutoHeal,
	}
	if err := encoder.Encode(ctri); err != nil {
		return fmt.Errorf("Error CreateMonitorFile: %w", err)
	}
	return nil
}

func aliveProcess(pid string) bool {
	pidAbs := filepath.Join("/", "proc", pid)
	_, err := os.Stat(pidAbs)
	return !os.IsNotExist(err)
}

func isAutoHeal(pid string) bool {
	h, _ := file.HomeDir("")
	f, _ := os.Open(filepath.Join(h, ".ops-controller", "instances", pid))
	decoder := json.NewDecoder(f)
	var i instance.ControllerInstance
	if err := decoder.Decode(&i); err != nil {
		fmt.Println("Decode Error while read AutoHeal Info")
		return false
	}
	if i.AutoHeal == true {
		return true
	} else {
		return false
	}
}
