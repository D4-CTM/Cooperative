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

func RenderTemplate(w http.ResponseWriter, path string, data map[string]string) {
  tmpl, err := template.ParseFiles(
    "./templates/layout.html",
    "./templates/"+ path +".html",
  )
  if err != nil {
    fmt.Printf("err.Error: %v\n", err.Error())
    http.Error(w, "Error loading template on " + path, http.StatusInternalServerError)
    return;
  }
  err = tmpl.Execute(w, data)
  if err != nil {
    fmt.Printf("err.Error: %v\n", err.Error())
    http.Error(w, "Error loading template on " + path, http.StatusInternalServerError)
    return;
  }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
  content := map[string]string{
    "title": "Login - ABCo-op",
    "stylePath": "login.css",
  }

  RenderTemplate(w, "login", content)  
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
  content := map[string]string{
    "title": "Register - ABCo-op",
    "stylePath": "register.css",
  }

  RenderTemplate(w, "register", content)
}

func main() {
	static := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", static))

  //mux.HandleFunc("/register-user/", registerUser)
	mux.HandleFunc("/verify-user/", verifyUser)
  mux.HandleFunc("/register", RegisterHandler)
  mux.HandleFunc("/login", LoginHandler)

	fmt.Println("Server started at:\nlocalhost:5312/")
	log.Fatal(http.ListenAndServe(":5412", mux))
}
