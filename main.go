package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func getRecourds(pageUrl string) {
	response := get("https://www.alqasidah.com/" + pageUrl)
	defer response.Body.Close()
	mainPageByte, err := ioutil.ReadAll(response.Body)
	errorCheck(&err)
	mainPageString := string(mainPageByte)
	//<u><b class='poemname'>هي لا تحبك أنت </b></u>
	name := mainPageString[strings.Index(mainPageString, "<u><b class='poemname'>")+23 : strings.Index(mainPageString, "</b></u>")-1]
	counter := ""
	for {

		index := strings.Index(mainPageString, ".mp3")
		if index < 0 {
			break
		}
		fmt.Printf("getting recourd: %v\n", name+counter)
		music := "https://www.alqasidah.com/" + mainPageString[index-17:index+4]
		filename := name + counter
		if len(filename) > 27 {
			filename = name[:26] + counter
		}
		file, err := os.Create(".\\records\\" + filename + ".mp3")
		errorCheck(&err)
		response := get(music)
		defer response.Body.Close()
		mainPageByte, err = ioutil.ReadAll(response.Body)
		errorCheck(&err)
		file.Write(mainPageByte)
		mainPageString = mainPageString[index+4:]
		counter += "1"
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
	fmt.Print("developed by ")
	fmt.Print(string("\033[32m"), "https://github.com/Kareem21227gg")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
func errorCheck(err *error) {
	if *err != nil {
		panic(*err)
	}
}

func get(url string) *http.Response {
	request, err := http.NewRequest("GET", url, nil)
	errorCheck(&err)
	request.Header.Add("Host", "alqasidah.com")
	response, err := http.DefaultClient.Do(request)
	errorCheck(&err)
	return response
}
