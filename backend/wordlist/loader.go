package wordlist

import (
	"log"

	"github.com/go-fsnotify/fsnotify"
	"github.com/huichen/sego"
)

var segmenter sego.Segmenter

// init
func InitLoader() {
	LoadDict()
	go func() {
		RunWatcher("data/sensitive.txt")
	}()
}

// load dict
func LoadDict() {
	dict := "data/sensitive.txt"
	segmenter.LoadDictionary(dict)
}

// api interface
func GetSegmenter() sego.Segmenter {
	return segmenter
}

// watch file change
func RunWatcher(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("file modified:" + event.Name)
					LoadDict()
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
