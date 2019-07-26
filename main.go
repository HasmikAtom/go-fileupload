package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File")
	r.ParseMultipartForm(10 << 20)

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

func routes() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8000", nil)
}

func main() {
	fmt.Println("Hello World")
	routes()
}
