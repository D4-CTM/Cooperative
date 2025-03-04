package main

import (
	"cooperative/backend"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func verifyUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\tVerifying user...")

	if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
		return
	}

	id := r.PostFormValue("code")
	pass := r.PostFormValue("password")

	user := backend.User{
		UserId:   id,
		Password: pass,
	}

	err := backend.Fetch(&user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
  
  fmt.Println("user_id:", user.UserId, "\tuser_pass:", user.Password)
	fmt.Println("User number:", user.UserNumber)

	fmt.Println("\tUser verify")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	pages := template.Must(template.ParseFiles("./templates/login.html"))
	pages.Execute(w, nil)
}

func main() {
	static := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", static))

	mux.HandleFunc("/verify-user/", verifyUser)
	mux.HandleFunc("/login", loginHandler)

	fmt.Println("starting connection with db")

	fmt.Println("Server started at:\nlocalhost:5312/")
	log.Fatal(http.ListenAndServe(":5412", mux))
}
