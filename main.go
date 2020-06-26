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

	vision "cloud.google.com/go/vision/apiv1"
)

const endpoint = "https://www.mhlw.go.jp/stf/seisakunitsuite/bunya/cocoa_00138.html"

func main() {
	body := request()
	data := parseImageDataByHTML(body)
	filePath := createImageFile(data)
	text := imageDetectDocumentText(filePath)

	defer os.Remove(filePath)

	fmt.Println(filePath)
	fmt.Println(text)
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

func imageDetectDocumentText(filePath string) string {
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
