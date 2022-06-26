package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type M map[string]interface{}

var tmpl *template.Template

type Info struct {
	Affiliation string
	Address     string
}

func (t Info) GetAffiliationDetailInfo() string {
	return "have 31 divisions"
}

type Person struct {
	Name    string
	Gender  string
	Hobbies []string
	Info    Info
}

type Superhero struct {
	Name    string
	Alias   string
	Friends []string
}

func (s Superhero) SayHello(from string, message string) string {
	return fmt.Sprintf("%s said: \"%s\"", from, message)
}

var funcMap = template.FuncMap{
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"avg": func(n ...int) int {
		var total = 0
		for _, each := range n {
			total += each
		}
		return total / len(n)
	},
}

func init() {
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/index", handlerIndex)
	http.HandleFunc("/hello", handlerHello)
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello again"))
	})
	http.HandleFunc("/about", handlerAbout)
	http.HandleFunc("/hero", handlerHero)
	http.HandleFunc("/custfunc", handlerCustFunc)
	http.HandleFunc("/form", routeIndexGet)
	http.HandleFunc("/process", routeSubmitPost)
	http.HandleFunc("/form-file", routeIndexGetFile)
	http.HandleFunc("/process-file", routeSubmitPostFile)
	http.HandleFunc("/ajax", handleAjax)
	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/getjson", actionIndex)
	http.HandleFunc("/multifile", handleMultifile)
	http.HandleFunc("/upload", handleUpload)

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

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	// var message = "Welcome"
	// // w.Write([]byte(message))
	// io.WriteString(w, message)

	// var filepath = path.Join("views", "index.html")
	// var tmpl, err = template.ParseFiles(filepath)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// var data = map[string]interface{}{
	// 	"title": "Learning Golang Web",
	// 	"name":  "Batman",
	// }

	// err = tmpl.Execute(w, data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	switch r.Method {
	case "POST":
		fmt.Fprint(w, "POST")
	case "GET":
		var person = Person{
			Name:    "Bruce Wayne",
			Gender:  "male",
			Hobbies: []string{"Reading Books", "Traveling", "Buying things"},
			Info:    Info{"Wayne Enterprises", "Gotham City"},
		}

		err := tmpl.ExecuteTemplate(w, "view", person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	var data = M{"name": "Batman"}
	err := tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	var message = "Hello world!"
	// w.Write([]byte(message))
	fmt.Fprint(w, message)
}

func handlerAbout(w http.ResponseWriter, r *http.Request) {
	var data = M{"name": "Batman"}
	err := tmpl.ExecuteTemplate(w, "about", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerHero(w http.ResponseWriter, r *http.Request) {
	var person = Superhero{
		Name:    "Bruce Wayne",
		Alias:   "Batman",
		Friends: []string{"Superman", "Flash", "Green Lantern"},
	}
	err := tmpl.ExecuteTemplate(w, "view-hero", person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerCustFunc(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Funcs(funcMap).ExecuteTemplate(w, "view-func", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routeIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("views/view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("views/view.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var name = r.FormValue("name")
		var message = r.Form.Get("message")

		var data = map[string]string{"name": name, "message": message}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeIndexGetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var tmpl = template.Must(template.New("form-file").ParseFiles("views/view.html"))
	var err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routeSubmitPostFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ...
	alias := r.FormValue("alias")

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

	// ...
	filename := handler.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}

	fileLocation := filepath.Join(dir, "files", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("done"))
}

func handleAjax(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/view-ajax.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := struct {
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Gender string `json:"gender"`
		}{}
		err := decoder.Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf(
			"Hllo, my name is %s. I'm %d years old %s",
			payload.Name,
			payload.Age,
			payload.Gender,
		)
		fmt.Fprint(w, message)
		return
	}

	http.Error(w, "Only accept POST request", http.StatusBadRequest)
}

func actionIndex(w http.ResponseWriter, r *http.Request) {
	data := []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{"Richard Grayson", 24},
		{"Jason Todd", 23},
		{"Tim Drake", 22},
		{"Damian Wayne", 21},
	}

	// jsonInBytes, err := json.Marshal(data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// fmt.Fprint(w, string(jsonInBytes))
	// // w.Write(jsonInBytes)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleMultifile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/view-multifile.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
		return
	}

	basePath, _ := os.Getwd()
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ...
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fileLocation := filepath.Join(basePath, "files", part.FileName())
		dst, err := os.Create(fileLocation)
		if dst != nil {
			defer dst.Close()
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(`all files uploaded`))
}
