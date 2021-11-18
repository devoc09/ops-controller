package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devoc09/ops-wrap/controller/instance"
)

type Monitor struct {
	targetDir string
	instance  *instance.Instance
	Msg       chan instance.Instance // to communicate with the controller
}

func New() (m *Monitor) {
	m = &Monitor{Msg: make(chan instance.Instance)}
	return
}

//process to be monitored every second
func (m *Monitor) Start(abs string) {
	fmt.Println(abs)
	m.targetDir = abs
	m.instance = instance.New(abs)
	m.Msg = make(chan instance.Instance)
	go func() {
		for {
			// fmt.Printf("Timer 2 second per Start........\n")
			ok := aliveProcess(filepath.Base(abs))
			if !ok {
				fmt.Printf("PID: %s Dead!!!\n", filepath.Base(abs))
				fmt.Printf("Send Channel Instance info\n")
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
	case healInstance := <-m.Msg:
		fmt.Printf("channel received......\n")
		fmt.Println(healInstance)
		err := exec.Command("ops", "instance", "create", "ops-hello", "-p", "8080").Start()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("instance heal Succeeded!!")
			return
		}
	}
}

func aliveProcess(pid string) bool {
	pidAbs := filepath.Join("/", "proc", pid)
	_, err := os.Stat(pidAbs)
	return !os.IsNotExist(err)
}
