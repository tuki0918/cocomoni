package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	vision "cloud.google.com/go/vision/apiv1"
)

const endpoint = "https://www.mhlw.go.jp/stf/seisakunitsuite/bunya/cocoa_00138.html"

type AppInfo struct {
	Version   string
	Date      string
	Downloads int
	Sentence  string
	Link      string
}

func main() {
	body := Request()
	data := ParseImageDataByHTML(body)
	filePath := CreateImageFile(data)
	text := DetectDocumentTextByImage(filePath)
	appInfo := ParseAppInfoByText(text)

	fmt.Println("[version]", appInfo.Version)
	fmt.Println("[date]", appInfo.Date)
	fmt.Println("[downloads]", appInfo.Downloads)
	fmt.Println("[sentence]", appInfo.Sentence)
	fmt.Println("[link]", appInfo.Link)

	defer os.Remove(filePath)
}

func Request() string {
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

func ParseImageDataByHTML(html string) string {
	r := regexp.MustCompile("<img src=\"data:image/png;base64,(.*?)\"[ /]*?>")
	data := r.FindStringSubmatch(html)[1]
	return data
}

func ParseAppInfoByText(text string) AppInfo {
	// r1 := regexp.MustCompile(`最新バージョンは「(\d+.\d+.\d+)」です。`)
	// version := r1.FindStringSubmatch(text)[1]
	r2 := regexp.MustCompile(`ダウンロード数は、(\d+月\d+日\d+:\d+)現在、合計で約(\d+)万件です。`)
	sentence := r2.FindStringSubmatch(text)[0]
	date := r2.FindStringSubmatch(text)[1]
	downloads := r2.FindStringSubmatch(text)[2]

	i, err := strconv.Atoi(downloads)
	if err != nil {
		log.Fatal(err)
	}

	appInfo := AppInfo{
		Version:   "x.x.x", // unsupported
		Date:      date,
		Downloads: i,
		Sentence:  sentence,
		Link:      endpoint,
	}
	return appInfo
}

func CreateImageFile(data string) string {
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

func DetectDocumentTextByImage(filePath string) string {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	return annotation.Text
}
