package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/form", handlerForm)
	http.HandleFunc("/process", routeSubmitPost)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(http.Dir("assets")))))

	var address = "localhost:9000"
	fmt.Printf("Server started at %s\n", address)
	http.ListenAndServe(address, nil)
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles(
		"views/index.html",
		"views/_header.html",
	))
	var data = map[string]interface{}{
		"title": "Index",
	}

	var err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerForm(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles(
		"views/form.html",
		"views/_header.html",
	))
	var data = map[string]interface{}{
		"title": "Form",
	}

	var err = tmpl.ExecuteTemplate(w, "form", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles(
			"views/result.html",
			"views/_header.html",
		))

		// Get the value of <input> form
		var name = r.FormValue("name")
		var city = r.FormValue("city")

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the value of <input file> form
		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uploadedFile, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		dir, err := os.Getwd()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := fmt.Sprintf("%s", handler.Filename)

		fileLocation := filepath.Join(dir, "upload", filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, uploadedFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = map[string]string{
			"name":     name,
			"city":     city,
			"fileLink": fileLocation,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Error(w, "", http.StatusBadRequest)
	}
}
