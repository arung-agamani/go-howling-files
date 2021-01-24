package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	isDev := os.Getenv("ISDEV") == "1"
	secureMiddleware := secure.New(secure.Options{
		ReferrerPolicy: "same-origin",
		AllowedHosts:   []string{"blog.howlingmoon.dev", "howlingmoon.dev"},
		IsDevelopment:  isDev,
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
	log.Println("Server started at port " + port)
	http.ListenAndServe(":"+port, router)
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
