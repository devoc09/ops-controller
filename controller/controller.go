package controller

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func newWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("NewWatcher() error: %w", err)
	}
	return watcher, nil
}

func Watch(target string) {
	watcher, _ := newWatcher()

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				// if event.Op&fsnotify.Write == fsnotify.Write {
				// 	log.Println("modified file:", event.Name)
				// }
				switch eventType := getEventType(event); eventType {
				// case "Write":
				case "Create":
					fmt.Printf("Create: %s\n", getFileName(event.Name))
					if err := writeCtrInstanceFile(event.Name); err != nil {
						fmt.Printf("Failed Create Controller Instance File: %s\n", getFileName(event.Name))
						fmt.Println(err)
						return
					}
				case "Remove":
					fmt.Printf("Remove: %s\n", getFileName(event.Name))
					if err := healInstance(event.Name); err != nil {
						fmt.Printf("Failed Heal Instance: %s\n", getFileName(event.Name))
						fmt.Println(err)
						return
					}
				// case "Rename":
				default:
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := watcher.Add(target)
	if err != nil {
		fmt.Println(err)
		return
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
