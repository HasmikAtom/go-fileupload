package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Receiving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprint(w, "Successfully uploaded File\n")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Downloading file...")

	rawURL := "https://d1ohg4ss876yi2.cloudfront.net/golang-resize-image/big.jpg"

	fileURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	path := fileURL.Path

	segments := strings.Split(path, "/")

	fileName := segments[2]

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	resp, err := check.Get(rawURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s with %v bytes downloaded", fileName, size)

}

func routes() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.ListenAndServe(":8000", nil)
}

func main() {
	fmt.Println("Hello World")
	routes()
}
