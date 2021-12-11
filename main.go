package main

import (
	"fmt"
	"html/template"
	"net/http"
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

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var name = r.FormValue("name")
		var city = r.FormValue("city")

		var data = map[string]string{
			"name": name,
			"city": city,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Error(w, "", http.StatusBadRequest)
	}
}
