package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func handel(w http.ResponseWriter, r *http.Request) {
	response := get("https://www.alqasidah.com/poet.php?poet=darwish")
	defer response.Body.Close()
	mainPageByte, err := ioutil.ReadAll(response.Body)
	errorCheck(&err)
	mainPageString := string(mainPageByte)
	mainPageString = mainPageString[65090:120860]

	titleList := strings.Split(mainPageString, "color:red")
	for i, title := range titleList {
		if i+1 == len(titleList) {
			titleList[i] = ""
			break
		}
		index := strings.LastIndex(title, "href")
		if index < 0 {
			titleList[i] = ""
			break
		}
		titleList[i] = title[index+6 : index+22]

	}
	fmt.Println("the urls list len: ", len(titleList))
	for _, pageUrl := range titleList {
		go getRecourds(pageUrl)
	}
	time.Sleep(time.Second * 20)
	w.Write([]byte(str))
}

var str string

func getRecourds(pageUrl string) {
	response := get("https://www.alqasidah.com/" + pageUrl)
	defer response.Body.Close()
	mainPageByte, err := ioutil.ReadAll(response.Body)
	errorCheck(&err)
	mainPageString := string(mainPageByte)
	var musicURLList = make([]string, 3)
	for i := 0; true; i++ {
		index := strings.Index(mainPageString, ".mp3")
		if index < 0 {
			break
		}
		musicURLList[i] = "<a href=\"https://www.alqasidah.com/" + mainPageString[index-17:index+4] + "\">music.mp3</a>\n"
		mainPageString = mainPageString[index+4:]
	}
	str += strings.Join(musicURLList, "")
}
func getPort() string {

	port := os.Getenv("PORT")
	if port == "" {
		return "3030"
	}
	return port
}
func main() {
	http.HandleFunc("/", handel)
	http.ListenAndServe(":"+getPort(), nil)

}
func errorCheck(err *error) {
	if *err != nil {
		panic(*err)
	}
}

func get(url string) *http.Response {
	response, err := http.Get(url)
	errorCheck(&err)
	return response
}
