package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
)

func main() {
	secureMiddleware := secure.New(secure.Options{
		ReferrerPolicy: "same-origin",
		AllowedHosts:   []string{"blog.howlingmoon.dev", "howlingmoon.dev"},
		IsDevelopment:  true,
	})
	router := mux.NewRouter()
	router.Use(secureMiddleware.Handler)
	router.Use(routeLogMiddleware)
	router.Handle("/api/ls", authMiddleware(http.HandlerFunc(fileListHandler)))
	router.Handle("/api/tree", authMiddleware(http.HandlerFunc(fileTreeHandler)))
	router.Handle("/upload", authMiddleware(routeLogMiddleware(http.HandlerFunc(uploadHandler)))).Methods("POST")

	// static file section
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./public/")})
	router.HandleFunc("/", index).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	log.Println("Server started at port 8080")
	http.ListenAndServe(":8080", router)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(path, "/") {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		return nil, err
	}
	return f, nil
}
