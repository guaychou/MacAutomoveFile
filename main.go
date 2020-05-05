package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"os/exec"
	"path/filepath"
	"regexp"
)

var extLocation = map[string]string{
	".jpg": "/Users/kevinchou/Pictures",
	".gif": "/Users/kevinchou/Pictures",
	".png": "/Users/kevinchou/Pictures",
	".txt": "/Users/kevinchou/Documents",
}
var patternWallpaper = regexp.MustCompile(`(?i)wallpaper`)
var patternSkripsi = regexp.MustCompile(`(?i)skripsi`)

var regexLocation = map[*regexp.Regexp]string{
	patternWallpaper: "<changeMe>",
	patternSkripsi: "<changeMe>",
}

func main() {
	// creates a new file watcher
	watcher,err := fsnotify.NewWatcher()
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Info("Created file: ", event.Name)
					dir,err:=getShouldDirectory(event.Name)
					if err!=nil{
						continue
					}
					cmd:=exec.Command("mv", event.Name, dir)
					err=cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
					log.Info("File moved")
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/Users/kevinchou/Desktop/test")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func getShouldDirectory(fileloc string) (string,error){
	extension:=filepath.Ext(fileloc)
	for pattern,value:=range(regexLocation){
		if pattern.MatchString(fileloc){
			return value,nil
		}
	}
	if extLocation[extension]==""{
		return "",errors.New("Extensions not listed")
	}
	return extLocation[extension],nil
}