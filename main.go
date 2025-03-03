package main

import (
	"fmt"
  "html/template"
	"log"
	"net/http"
  "cooperative/backend"
)

func verifyUser(w http.ResponseWriter, r *http.Request) {
  fmt.Println("\n\tVerifying user...")

  if (r.Header.Get("HX-Request") == "") {
    fmt.Println("\tWASN'T A HX-Request!")
    return
  }

  code := r.PostFormValue("code")
  pass := r.PostFormValue("password")

  fmt.Println("User's code:", code, "\nUsers's pass:" , pass)

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
  backend.TestConnection()

	fmt.Println("Server started at:\nlocalhost:5312/")
	log.Fatal(http.ListenAndServe(":5412", mux))
}
