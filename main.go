package main

import (
	"fmt"
	"net/http"
	"path"
	"text/template"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/index", handlerIndex)
	http.HandleFunc("/hello", handlerHello)
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello again"))
	})

	var address = ":8080"
	fmt.Printf("server started at %s\n", address)

	// err := http.ListenAndServe(address, nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// using http.Server
	var server *http.Server = new(http.Server)
	server.Addr = address
	// server.ReadTimeout = time.Second * 10
	// server.WriteTimeout = time.Second * 10
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	// var message = "Welcome"
	// // w.Write([]byte(message))
	// io.WriteString(w, message)

	var filepath = path.Join("views", "index.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"title": "Learning Golang Web",
		"name":  "Batman",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	var message = "Hello world!"
	// w.Write([]byte(message))
	fmt.Fprint(w, message)
}
