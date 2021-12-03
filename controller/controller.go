package controller

import (
	"fmt"
	"log"

	"github.com/devoc09/ops-wrap/internal/monitor"
	"github.com/fsnotify/fsnotify"
)

func newWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("NewWatcher() error: %w", err)
	}
	return watcher, nil
}

func Watch(targets []string) {
	watcher, _ := newWatcher()

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					done <- true
					return
				}
				// monCh := make(chan *instance.Instance)

				switch eventType := getEventType(event); eventType {
				case "Create":
					m := monitor.New(event.Name)
					fmt.Printf("Create: %s\n", getFileName(event.Name))
					go func() {
						m.Start(event.Name)
					}()
					go func() {
						if err := m.CreateMonitorFile(event.Name); err != nil {
							return
						}
					}()
				case "Remove":

				}
				// // process branching to a message from monitor
				// select {
				// case monMsg := <-mon.Msg:
				// 	fmt.Printf("channel received......\n")
				// 	fmt.Println(monMsg)
				// 	err := exec.Command("ops", "instance", "create", "ops-hello", "-p", "8080").Start()
				// 	if err != nil {
				// 		fmt.Println(err)
				// 	} else {
				// 		fmt.Println("instance heal Succeeded!!")
				// 	}
				// }
			case err, ok := <-watcher.Errors:
				if !ok {
					done <- true
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// add watching directory
	for _, t := range targets {
		err := watcher.Add(t)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	<-done
}

func getEventType(event fsnotify.Event) string {
	if event.Op&fsnotify.Write == fsnotify.Write {
		return "Write"
	} else if event.Op&fsnotify.Create == fsnotify.Create {
		return "Create"
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		return "Remove"
	} else if event.Op&fsnotify.Rename == fsnotify.Rename {
		return "Rename"
	} else {
		return "Chmod"
	}
}
