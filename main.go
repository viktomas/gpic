package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed assets/index.html
var indexTemplate string

//go:embed assets/style.css
var css string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: picsort <folder>")
		return
	}

	rootFolder := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").Parse(indexTemplate))
		imageFiles := getImageFiles(rootFolder)

		err := tmpl.Execute(w, imageFiles)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(rootFolder))))

	http.HandleFunc("/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		io.WriteString(w, css)
	})

	addr := "localhost:8080"
	fmt.Printf("Starting server at http://%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func getImageFiles(folderPath string) []string {
	var imageFiles []string
	_ = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(path) {
			imageFiles = append(imageFiles, filepath.Base(path))
		}
		return nil
	})
	return imageFiles
}

func isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
