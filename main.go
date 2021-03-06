package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatalf("Expecting at least 2 arguments got %d", len(os.Args)-1)
	}

	watcher, err := fsnotify.NewWatcher()
	execLoop := NewExecLoop(os.Args[2], os.Args[3:])
	done := make(chan bool)

	if err != nil {
		log.Fatalf("Failed to create new watcher %s", err)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("\x1b[0;32mChanged %s\x1b[0m", event.String())

				if event.Op&fsnotify.Create == fsnotify.Create {
					stats, err := os.Stat(event.Name)

					if err != nil {
						log.Printf("Failed to stat changed thing %s", err)
					} else if stats.IsDir() {
						log.Printf("\x1b[0;34mWill add dir to watch `%s`\x1b[0m", event.Name)
						watcher.Add(event.Name)
					}
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("\x1b[0;34mWill remove watch path `%s`\x1b[0m", event.Name)
					watcher.Remove(event.Name)
				}

				execLoop.Exec()

			case err := <-watcher.Errors:
				log.Printf("Got failure %s", err)
			}
		}
	}()

	log.Printf("Will watch %s with command %s", os.Args[1], strings.Join(os.Args[2:], " "))

	err = watcher.Add(os.Args[1])
	log.Printf("\x1b[0;34mAdding path `%s`\x1b[0m", os.Args[1])

	if err != nil {
		log.Fatalf("Failed to add path to watch %s", err)
	}

	addRecursively(watcher, os.Args[1])

	<-done
}

func addRecursively(watcher *fsnotify.Watcher, path string) {
	paths := []string{path}

	for i := 0; i < len(paths); i++ {
		path = paths[i]
		stats, err := os.Stat(path)

		if err == nil && stats.IsDir() {
			file, err := os.Open(path)
			if err == nil {
				subs, _ := file.Readdir(0)

				for _, sub := range subs {
					if sub.IsDir() {
						fullPath := filepath.Join(path, sub.Name())
						paths = append(paths, fullPath)
						watcher.Add(fullPath)
						log.Printf("\x1b[0;34mAdding path `%s`\x1b[0m", fullPath)
					}
				}
			}
		}
	}
}
