package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type User_register struct {
	Name, Email, Password string
	//token string
}

var m = map[string]string{
	"Name":     "Name",
	"Email":    "Example@.com",
	"Password": "Password",
}

func homepage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/homepage.html")
	tmpl.Execute(w, nil)

}
func register(w http.ResponseWriter, r *http.Request) {
	reg, _ := template.ParseFiles("templates/register.html")
	reg.Execute(w, nil)
}

func authorize(w http.ResponseWriter, r *http.Request) {
	reg, _ := template.ParseFiles("templates/authorize.html")
	reg.Execute(w, m)
	if val, ok := m["Email"]; ok {
		fmt.Println(val)
	}
	reg.Execute(w, m)
	fmt.Println(m)

}
func account(w http.ResponseWriter, r *http.Request) {
	reg, _ := template.ParseFiles("templates/account.html")

	m["Name"] = r.FormValue("Name")
	m["Email"] = r.FormValue("Email")
	m["Password"] = r.FormValue("Password")
	/*	if val, ok := m["Email"]; ok{
	    fmt.Println(val)
		}*/
	reg.Execute(w, m)
	fmt.Println(m)
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/homepage", homepage)
	rtr.HandleFunc("/register", register).Methods("GET")
	rtr.HandleFunc("/authorize", authorize).Methods("GET")
	rtr.HandleFunc("/account", account)
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
