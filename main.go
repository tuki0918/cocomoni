package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const endpoint = "https://www.mhlw.go.jp/stf/seisakunitsuite/bunya/cocoa_00138.html"

func main() {
	body := request()
	data := parseImageDataByHTML(body)
	filePath := createImageFile(data)

	defer os.Remove(filePath)

	fmt.Println(filePath)
}

func request() string {
	response, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func parseImageDataByHTML(html string) string {
	r := regexp.MustCompile("<img src=\"data:image/png;base64,(.*?)\"[ /]*?>")
	data := r.FindStringSubmatch(html)[1]
	return data
}

func createImageFile(data string) string {
	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatal(err)
	}
	tmpFile, err := ioutil.TempFile("", "*.png")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Write(dec); err != nil {
		tmpFile.Close()
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpFile.Name()
}
