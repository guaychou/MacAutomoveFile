package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	patternWallpaper = regexp.MustCompile(`(?i)wallpaper`)
	patternSkripsi = regexp.MustCompile(`(?i)skripsi`)
	regexPicture = regexp.MustCompile(`jpeg|jpg|png|gif|tif`)
	regexMusic = regexp.MustCompile(`mp3|m4a|wav`)
	regexMovie=regexp.MustCompile(`mp4|mov|mkv`)
	regexDocument=regexp.MustCompile(`pdf|doc|docx|odt|odp|pptx|ppt`)
	regexChromeDlExt=regexp.MustCompile(`crdownload`)
	regexCompressedDoc=regexp.MustCompile(`zip|tar.gz`)
)

var (
	regexLocation = map[*regexp.Regexp]string{
		patternWallpaper: "/Users/kevinchou/Library/Mobile Documents/com~apple~CloudDocs/Wallpaper",
		patternSkripsi: "/Users/kevinchou/Library/Mobile Documents/com~apple~CloudDocs/Documents/Skripsi",
	}
	regexExtension = map[*regexp.Regexp]string{
		regexPicture: "/Users/kevinchou/Pictures",
		regexMovie: "/Users/kevinchou/Movies",
		regexMusic: "/Users/kevinchou/Music",
		regexDocument: "/Users/kevinchou/Documents",
	}
    watcher *fsnotify.Watcher

)

func main() {
	// creates a new file watcher
	watcher,_=fsnotify.NewWatcher()
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

	if err := filepath.Walk("/Users/kevinchou/Desktop/test", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}

func getShouldDirectory(fileloc string) (string,error){
	extension:=filepath.Ext(fileloc)
	if regexChromeDlExt.MatchString(extension){
		return "",errors.New("Chrome's download file, continuing. . . ")
	}
	if regexCompressedDoc.MatchString(extension){
		uncompress(fileloc,extension)
	}
	for pattern,value:=range(regexLocation){
		if pattern.MatchString(fileloc){
			return value,nil
		}
	}
	for pattern,value:=range(regexExtension){
		if pattern.MatchString(extension){
			return value,nil
		}
	}
	return "", errors.New("cannot found directory for this extensions")
}

func uncompress(fileloc string,extension string){
	if extension==".tar.gz"{
		cmd:=exec.Command("tar","-xzf", fileloc ,"-C", filepath.Dir(fileloc))
		err:=cmd.Run()
		if err!=nil{
			log.Fatal(err)
		}
	}else{
		cmd:=exec.Command("unzip", fileloc, "-d",filepath.Dir(fileloc))
		err:=cmd.Run()
		if err!=nil{
			log.Fatal(err)
		}
	}
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}