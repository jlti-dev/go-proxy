package main

import (
	"fmt"
	"os"

	"net/http"
	"crypto/tls"

	"io"
	"path/filepath"
	"strings"

	"log"
	"time"

)

func onlineCache(filename string) {
	//get file info
	fileInfo, err := os.Stat(fmt.Sprintf("%s%s", STATIC_DIR, filename))
	//handle error
	if err != nil {
		log.Println(err)
		downloadFile(filename)
	}else{
		timeDiff := time.Now().Sub(fileInfo.ModTime())
		
		if timeDiff.Minutes() > float64(heatMins) {
			log.Printf("%s was modified at %s, that was before %f minutes", fileInfo.Name(), fileInfo.ModTime().Format(time.UnixDate), timeDiff.Minutes())
			offlineCache(filename)
		}else{
			log.Printf("%s is still valid in heatMins", filename)
		}
	}
}
func offlineCache(filename string){
	log.Println("Offline Cache is asynchron")
	go downloadFile(filename)
}
func downloadFile(urlsuffix string){
	url := fmt.Sprintf("%s%s",baseUrl, urlsuffix)
	
	log.Println("Start Downloading %s", url)

	if err := os.MkdirAll(filepath.Dir(fmt.Sprintf("%s%s", STATIC_DIR, urlsuffix)), 0660); err != nil {
		log.Println(err)
		return
	}
	// Create blank file
        //Parameter hier abschneiden!
	fileName := strings.Split(fmt.Sprintf("%s%s", STATIC_DIR, urlsuffix), "?")[0]
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
		Transport: tr,
	}
	// Put content on file
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	defer file.Close()
	
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		size, err := io.Copy(file, resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Downloaded %s with size %d into %s", url, size, fileName)
	}else {
		log.Printf("Statuscode der Anfrage %s: %d", url, resp.StatusCode)
	}
	

}

