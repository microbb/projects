package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Dir struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
	Dirs  []*Dir   `json:"dirs"`
}

var path string

func main() {

	fmt.Print("Введите путь до директории: ")
	fmt.Scanln(&path)

	http.HandleFunc("/download", HandleDownload)
	http.HandleFunc("/structure", HandleGetStructure)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("Server run fail")
	}

}

func HandleGetStructure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := generateData(path)

	res, _ := json.MarshalIndent(data, "", "\t")

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")

	ff, err := os.Open(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Unable to open file ")
		return
	}

	ff.Close()

	// force a download with the content- disposition field
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(name))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, name)
}

// Рекурсивный проход по директориям и запись их структурно в json
func generateData(root string) *Dir {

	mainDir := Dir{
		Name: root,
	}

	entities, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}

	if len(entities) == 0 {
		return &mainDir
	}

	for _, entity := range entities {
		if entity.IsDir() {
			subDir := generateData(filepath.Join(root, entity.Name()))
			mainDir.Dirs = append(mainDir.Dirs, subDir)
		} else {
			mainDir.Files = append(mainDir.Files, filepath.Join(root, entity.Name()))
		}
	}

	return &mainDir
}
