package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
func writeJSON(f string, t Tickets) {
	file, _ := json.MarshalIndent(t, "", " ")
	_ = ioutil.WriteFile(f, file, 0644)

}

//Go application entrypoint
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tickets := readJSON("tickets.json")
		templates := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
		if err := templates.ExecuteTemplate(w, "layout.html", tickets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	r.HandleFunc("/tickets/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		tickets := readJSON("tickets.json")
		var ticket Ticket
		for _, t := range tickets.Tickets {
			if t.ID == id {
				ticket = t
			}
		}
		if ticket != (Ticket{}) {
			templates := template.Must(template.ParseFiles("templates/layout.html", "templates/ticket.html"))
			if err := templates.ExecuteTemplate(w, "layout.html", ticket); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "The requested resource could not be found", http.StatusNotFound)
		}
	})
	r.HandleFunc("/tickets/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		tickets := readJSON("tickets.json")
		var newTickets Tickets
		for i, t := range tickets.Tickets {
			if t.ID == id {
				newTickets.Tickets = append(tickets.Tickets[:i], tickets.Tickets[i+1:]...)
				writeJSON("tickets.json", newTickets)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	r.HandleFunc("/tickets/close/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		tickets := readJSON("tickets.json")
		for i, t := range tickets.Tickets {
			if t.ID == id {
				tickets.Tickets[i].Status = "closed"
				writeJSON("tickets.json", tickets)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	r.HandleFunc("/tickets/open/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		tickets := readJSON("tickets.json")
		for i, t := range tickets.Tickets {
			if t.ID == id {
				tickets.Tickets[i].Status = "open"
				writeJSON("tickets.json", tickets)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	//  Handle static content
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Println("Server started on port 8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}
