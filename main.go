package main

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
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
		imageFiles, err := getImageFiles(rootFolder)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		err = tmpl.Execute(w, imageFiles)
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

	http.HandleFunc("/delete-similar", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing form %v", err.Error()), 500)
			return
		}

		toDeleteFolder := "to-delete"
		err = os.MkdirAll(filepath.Join(rootFolder, toDeleteFolder), 0755)
		if err != nil {
			http.Error(w, fmt.Sprintf("error creating delete folder: %v", err.Error()), 500)
			log.Printf("error creating delete folder: %v", err.Error())
			return
		}
		similarPictures := []string{}
		similarIndex := 0
		for {
			key := fmt.Sprintf("similar[%d]", similarIndex)
			similarPicture := r.Form.Get(key)
			if similarPicture == "" {
				break
			}
			similarPictures = append(similarPictures, similarPicture)
			similarIndex++
		}
		survivors := map[string]struct{}{}
		for image := range r.Form {
			survivors[image] = struct{}{}
		}
		for _, image := range similarPictures {
			if _, ok := survivors[image]; ok {
				continue
			}
			err = os.Rename(filepath.Join(rootFolder, image), filepath.Join(rootFolder, toDeleteFolder, image))
			if err != nil {
				http.Error(w, fmt.Sprintf("error moving image: %v", err.Error()), 500)
				log.Printf("error moving image: %v", err.Error())
				return
			}
		}
		http.RedirectHandler("/similar", 301).ServeHTTP(w, r)

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

func getImageFiles(folderPath string) ([]string, error) {
	var imageFiles []string
	allFiles, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}
	for _, file := range allFiles {

		if !file.IsDir() && isImageFile(file.Name()) {
			imageFiles = append(imageFiles, file.Name())
		}
	}
	return imageFiles, nil
}

func isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
