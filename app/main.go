package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

var dbUser, dbPassword, dbHost string

func setVariables() {
	dbUser = os.Getenv("db_user")
	dbPassword = os.Getenv("db_password")
	dbHost = os.Getenv("db_host")
}
func readData() []Ticket {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":3306)/tickets")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	results, err := db.Query("SELECT id, t_title, t_desc, t_time, t_status FROM tickets")
	if err != nil {
		panic(err.Error())
	}
	var tickets []Ticket
	for results.Next() {
		var t Ticket
		// for each row, scan the result into our tag composite object
		err = results.Scan(&t.ID, &t.Title, &t.Description, &t.Time, &t.Status)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		tickets = append(tickets, t)
	}

	return tickets
}
func deleteTicket(id string) {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":3306)/tickets")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	delForm, err := db.Prepare("DELETE FROM tickets WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)
}

func closeTicket(id string) {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":3306)/tickets")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	closeForm, err := db.Prepare("UPDATE tickets SET t_status = 'closed' WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	closeForm.Exec(id)
}

func openTicket(id string) {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":3306)/tickets")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	closeForm, err := db.Prepare("UPDATE tickets SET t_status = 'open' WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	closeForm.Exec(id)
}

//Go application entrypoint
func main() {
	setVariables()
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tickets := readData()
		templates := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
		if err := templates.ExecuteTemplate(w, "layout.html", tickets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	r.HandleFunc("/tickets/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		tickets := readData()
		var ticket Ticket
		for _, t := range tickets {
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
		deleteTicket(id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	r.HandleFunc("/tickets/close/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		closeTicket(id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	r.HandleFunc("/tickets/open/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		openTicket(id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	//  Handle static content
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Println("Server started on port 8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}
