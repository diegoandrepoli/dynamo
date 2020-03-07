package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

/**
 * Path of configuration files
 */
const ConfigPath = "files/"

/**
 * Configuration file extension
 */
const FileExtension = ".csv"

/**
 * Main application
 */
func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/file/{name}", create).Methods("POST")
	router.HandleFunc("/file/{name}", append).Methods("PUT")
	router.HandleFunc("/file/{name}", read).Methods("GET")
	router.HandleFunc("/file", list).Methods("GET")

	fmt.Println("Running config server at port 8080!")
	log.Fatal(http.ListenAndServe(":8080", router))
}

/**
 * Index application response
 */
func index(w http.ResponseWriter, _ *http.Request){
	fmt.Fprintln(w, "Hello, welcome to Dynamo!")
}

/**
 * Path of file
 */
func path(name string) string {
	return fmt.Sprintf("%s%s%s", ConfigPath, name, FileExtension)
}

/**
 * Create file
 */
func create(w http.ResponseWriter, r *http.Request)  {
	s := mux.Vars(r)["name"]
	if s == "" {
		http.Error(w, "Invalid file name", 500)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}

	err = ioutil.WriteFile(path(mux.Vars(r)["name"]), b, 0644)
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}
}

/**
 * Append file
 */
func append(w http.ResponseWriter, r *http.Request)  {
	s := mux.Vars(r)["name"]
	if s == "" {
		http.Error(w, "Invalid file name", 500)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}

	file, err := os.OpenFile(path(mux.Vars(r)["name"]), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed opening file: %s", err), 500)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\r\n%s", b))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to writing file: %s", err), 500)
		return
	}
}

/**
 * Read file
 */
func read(w http.ResponseWriter, r *http.Request)  {
	b, err := ioutil.ReadFile(path(mux.Vars(r)["name"]))
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}

	fmt.Fprintln(w, string(b))
}

/**
 * List files
 */
func list(w http.ResponseWriter, r *http.Request){
	files, err := ioutil.ReadDir(ConfigPath)
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}

	for _, f := range files {
		fmt.Fprintln(w, strings.Replace(f.Name(), FileExtension, "", -1))
	}
}


