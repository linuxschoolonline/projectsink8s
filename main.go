package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Ticket a struct that holds information to be displayed in our HTML file
type Ticket struct {
	ID          string
	Time        string
	Title       string
	Description string
	Status      string
}

// Tickets holds all the Tickets
type Tickets struct {
	Tickets []Ticket `json:"tickets"`
}

func readJSON(f string) Tickets {
	jsonFile, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	var tickets Tickets
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &tickets)
	return tickets
}

//Go application entrypoint
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tickets := readJSON("tickets.json")
		templates := template.Must(template.ParseFiles("templates/home.html"))
		if err := templates.ExecuteTemplate(w, "home.html", tickets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println(http.ListenAndServe(":8080", nil))
}
