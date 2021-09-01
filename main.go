package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dhowden/tag"
)

func getRecourds(pageUrl string) {
	response := get("https://www.alqasidah.com/" + pageUrl)
	defer response.Body.Close()
	mainPageByte, err := ioutil.ReadAll(response.Body)
	errorCheck(&err)
	mainPageString := string(mainPageByte)
	name := "blh.mp3"
	for {
		index := strings.Index(mainPageString, ".mp3")
		if index < 0 {
			break
		}
		music := "https://www.alqasidah.com/" + mainPageString[index-17:index+4]
		file, err := os.Create(".\\records\\" + name)
		errorCheck(&err)
		response := get(music)
		defer response.Body.Close()
		mainPageByte, err = ioutil.ReadAll(response.Body)
		errorCheck(&err)
		file.Write(mainPageByte)
		tag.ReadFrom(file)
		mainPageString = mainPageString[index+4:]
	}
}
func main() {
	os.Mkdir("records", fs.ModeDir)
	response := get("https://www.alqasidah.com/poet.php?poet=darwish")
	defer response.Body.Close()
	mainPageByte, err := ioutil.ReadAll(response.Body)
	errorCheck(&err)
	mainPageString := string(mainPageByte)
	mainPageString = mainPageString[65090:120860]
	//all poems with record has "color:red" in thir tags
	titleList := strings.Split(mainPageString, "color:red")
	titleList[len(titleList)-1] = ""
	for i, title := range titleList {
		index := strings.LastIndex(title, "href")
		if index < 0 {
			i--
			continue
		}
		titleList[i] = title[index+6 : index+22]
	}
	fmt.Println("the urls list length: ", len(titleList))
	for _, pageUrl := range titleList {
		getRecourds(pageUrl)
	}

}
func errorCheck(err *error) {
	if *err != nil {
		panic(*err)
	}
}

func get(url string) *http.Response {
	//this is importanat to avoid http error 403
	time.Sleep(time.Second / 2)
	request, err := http.NewRequest("GET", url, nil)
	errorCheck(&err)
	request.Header.Add("Host", "alqasidah.com")
	response, err := http.DefaultClient.Do(request)
	errorCheck(&err)
	return response
}
