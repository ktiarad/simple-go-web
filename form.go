package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

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

		// Rename filename
		filename := fmt.Sprintf("%s", handler.Filename)
		newFilename := fmt.Sprintf("%s%s", renameFile(filename), filepath.Ext(handler.Filename))

		fileLocation := filepath.Join(dir, "upload", newFilename)
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
			"fileName": newFilename,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Error(w, "", http.StatusBadRequest)
	}
}

func renameFile(filename string) string {
	var salt = fmt.Sprintf("%d", time.Now().UnixNano())
	var salted = fmt.Sprintf("%s%s", filename, salt)

	var sha = sha1.New()
	sha.Write([]byte(salted))
	var newFilename = sha.Sum(nil)

	return fmt.Sprintf("%x", newFilename)
}

func handlerDownload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.FormValue("file")
	basePath, _ := os.Getwd()
	fileLocation := filepath.Join(basePath, "upload", path)

	f, err := os.Open(fileLocation)
	if f != nil {
		defer f.Close()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", f.Name())
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
