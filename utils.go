package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/h2non/bimg"
)

func fileWriter(filepath string, data []byte) bool {
	log.Println(len(data))
	file, err := os.Create(filepath)
	if err != nil {
		return false
	}
	if _, err := file.Write(data); err != nil {
		return false
	}
	defer file.Close()
	return true
}

func convertToWebp(f multipart.File) ([]byte, bool) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		log.Println("copy: ", err)
		return nil, false
	}
	imageBuffer := buf.Bytes()
	newImage, err := bimg.NewImage(imageBuffer).Convert(bimg.WEBP)
	if err != nil {
		log.Println("convert: ", err, len(buf.Bytes()))
		return buf.Bytes(), false
	}
	if bimg.NewImage(newImage).Type() == "webp" {
		return newImage, true
	}
	return buf.Bytes(), false
}

func constructDirTree(root DirTree) DirTree {
	if root.IsFile {
		leaf := DirTree{
			Name:       root.Name,
			IsFile:     true,
			Child:      []DirTree{},
			CurrentDir: root.CurrentDir,
		}
		root.Child = []DirTree{}
		return leaf
	} // is a directory, contains childs
	children := []DirTree{}
	files, err := ioutil.ReadDir(root.CurrentDir)
	if err != nil {
		// error on this directory
		leaf := DirTree{
			Name:   root.Name,
			IsFile: false,
			Child:  children,
		}
		return leaf
	}
	for file := range files {
		isDir := files[file].IsDir()
		child := DirTree{
			Name:       files[file].Name(),
			CurrentDir: root.CurrentDir + "/" + files[file].Name(),
			Child:      []DirTree{},
			IsFile:     !isDir,
		}
		children = append(children, constructDirTree(child))
	}
	this := DirTree{
		Name:       root.Name,
		IsFile:     false,
		CurrentDir: root.CurrentDir,
		Child:      children,
	}
	return this
}

func compareTime(a time.Time, b time.Time) int {
	if a.After(b) {
		return 1
	} else if a.Equal(b) {
		return 0
	} else {
		return -1
	}
}

func setJSONHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	return w
}

func parseDirQuery(q url.Values) (bool, int, int) {
	month := q.Get("month")
	day := q.Get("day")
	intMonth, err := strconv.Atoi(month)
	intDay, err2 := strconv.Atoi(day)
	if err != nil || err2 != nil {
		return false, 0, 0
	}
	return true, intMonth, intDay
}
