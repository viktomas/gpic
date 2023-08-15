package main

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

//go:embed templates/similar.html
var simlarTemplate string

//go:embed templates/compare-similar.html
var compareSimilarTemplate string

//go:embed assets
var assets embed.FS

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: picsort <folder>")
		return
	}

	rootFolder := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			io.WriteString(w, "Page not found.")
			return
		}
		http.RedirectHandler("/similar", 301).ServeHTTP(w, r)

	})

	http.HandleFunc("/similar", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").Parse(simlarTemplate))
		imageFiles := getImageFiles(rootFolder)

		err := tmpl.Execute(w, imageFiles)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	http.HandleFunc("/compare-similar", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("compare-similar").Parse(compareSimilarTemplate))
		imageFiles := []string{}

		for k := range r.URL.Query() {

			imageFiles = append(imageFiles, k)
		}
		sort.Strings(imageFiles)
		grid := 9
		if len(imageFiles) <= 2 {
			grid = 2
		} else if len(imageFiles) <= 4 {
			grid = 4
		}
		err := tmpl.Execute(w, map[string]any{
			"grid":     grid,
			"pictures": imageFiles,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(rootFolder))))

	fs := http.FileServer(http.FS(assets))
	http.Handle("/assets/", fs)

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
