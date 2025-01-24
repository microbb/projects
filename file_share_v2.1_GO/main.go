package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ResData struct {
	Name  string `json:"name"`
	Days  int    `json:"days"`
	Files []File
}

type File struct {
	Name string `json:"fileName"`
	Url  string `json:"url"`
}

func main() {
	wg := sync.WaitGroup{}

	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		data := ResData{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range data.Files {
			wg.Add(1)

			go DownloadFile(file.Name, data.Name, file.Url, &wg)
		}

	})

	wg.Wait()

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func DownloadFile(fileName, archiveName, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Не удалось получить директорию в которой запущен сервис")
	}

	x := strings.Split(fileName, ".")[0]

	pathToFile := fmt.Sprintf("%s/tmp/%s/%s", dir, x, fileName)
	//fmt.Println(pathToFile)

	_ = os.MkdirAll(dir+"/tmp/"+x, os.ModePerm)

	file, err := os.Create(pathToFile)
	if err != nil {
		fmt.Println("Не удалось создать файл")
	}

	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Не удалось сделать запрос")
	}

	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Ошибка при записи данных в файл")
	}

	fmt.Println("файл скачен")
	ArchiveFiles(fileName, archiveName, "test")

}

func ArchiveFiles(fileName, nameArchive, passwd string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/tmp/%s", dir, strings.Split(fileName, ".")[0])

	cmd := exec.Command("7z", "a", "-tzip", "-mx1", "-sdel", "-p"+passwd, nameArchive)
	cmd.Dir = path
	cmd.Run()

	fmt.Println("Файл заархивирован")
}
