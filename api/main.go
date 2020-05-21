package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
func main() {
	setVariables()
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tickets := readData()
		js, err := json.Marshal(tickets)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	log.Println("Server started on port 8000")
	fmt.Println(http.ListenAndServe(":8000", r))
}
