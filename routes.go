package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/h2non/bimg"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := Message{
		Status:  "success",
		Message: "Hello World",
	}
	json.NewEncoder(w).Encode(res)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	userAuth, err := getUserAuth(key)
	if err != nil {
		res := Message{
			Status:  "failed",
			Message: "There is some weird error",
		}
		log.Println(err)
		json.NewEncoder(w).Encode(res)
		return
	}
	res := UserAuthMessage{
		Status: "success",
		Data:   userAuth,
	}

	json.NewEncoder(w).Encode(res)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var isAnyError = false
	var res Message
	var writtenCount int = 0
	r.ParseMultipartForm(32 << 20)
	fhs, _ := r.MultipartForm.File["file"]

	if len(fhs) == 0 {
		res = Message{
			Status:  "failed",
			Message: "no file received",
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	year, month, day := time.Now().Date()
	currentDir := fmt.Sprintf("%d-%d-%d", year, int(month), int(day))
	dirString := filepath.Join("./", "public", "blog", currentDir)

	// make directory if not exists
	if _, err := os.Stat(filepath.Join("./", "public", "blog")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join("./", "public", "blog"), os.ModeDir)
	}
	if _, err := os.Stat(dirString); os.IsNotExist(err) {
		os.Mkdir(dirString, os.ModeDir)
	}

	for _, fh := range fhs {
		fileName := strings.TrimSuffix(fh.Filename, filepath.Ext(fh.Filename))
		file, err2 := fh.Open()
		if err2 != nil {
			fmt.Println(err2)
			isAnyError = true
		}
		buf := make([]byte, 512)
		if _, err := file.Read(buf); err != nil {
			log.Println(err)
			continue
		}
		file.Seek(0, io.SeekStart)
		if strings.HasPrefix(http.DetectContentType(buf), "image") {
			isImage := true
			newFileBuf, stat := convertToWebp(file)
			if stat == false {
				if isImage {
					prefix := fmt.Sprintf("%d-%s", time.Now().Unix(), fh.Filename)
					fhr, err := fh.Open()
					if err != nil {
						continue
					}
					defer fhr.Close()
					dst, err := os.Create(filepath.Join(dirString, prefix))
					defer dst.Close()
					if _, err := io.Copy(dst, fhr); err != nil {
						isAnyError = true
					} else {
						writtenCount++
					}

				} else {
					isAnyError = true
				}
			} else {
				prefix := fmt.Sprintf("%d-%s.webp", time.Now().Unix(), fileName)
				bimg.Write(filepath.Join(dirString, prefix), newFileBuf)
				defer file.Close()
				writtenCount++
			}
		}

	}

	if isAnyError != true {
		res = Message{
			Status:  "success",
			Message: fmt.Sprintf("Files uploaded successfully, %d files written", writtenCount),
		}
	} else {
		res = Message{
			Status:  "failed",
			Message: "There is file error",
		}
	}
	json.NewEncoder(w).Encode(res)
}

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	isValid, month, day := parseDirQuery(query)
	if !isValid {
		res := Message{
			Status:  "failed",
			Message: "Invalid request",
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}
	year, _, _ := time.Now().Date()
	imgDirName := fmt.Sprintf("%d-%d-%d", year, month, day)
	log.Println(imgDirName)
	files, err := ioutil.ReadDir(filepath.Join("./", "public", "blog", imgDirName))
	if err != nil {
		res := Message{
			Status:  "failed",
			Message: "Internal Server Error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(res)
		log.Println(err)
		return
	}
	var dirNameArr []string
	for file := range files {
		log.Println(files[file].Name())
		if !files[file].IsDir() {
			dirName := fmt.Sprintf("public/blog/%s/%s", imgDirName, files[file].Name())
			dirNameArr = append(dirNameArr, dirName)
		}
	}

	res := DirListMessage{
		Status: "success",
		Data:   dirNameArr,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func fileTreeHandler(w http.ResponseWriter, r *http.Request) {
	currentDir := filepath.Join("./", "public")
	root := DirTree{
		Name:       "public",
		IsFile:     false,
		CurrentDir: currentDir,
		Child:      []DirTree{},
	}
	res := DirTreeMessage{
		Status: "maybe",
		Data:   constructDirTree(root),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
