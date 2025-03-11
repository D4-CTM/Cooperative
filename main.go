package main

import (
	"cooperative/backend"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func registerUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Starting registration user...")

    if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
        return ;
    }

}

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
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`<p>There was an issue finding the user, please verify<br>the password or the user.</p>`))
		return
	}

	fmt.Println("user_id:", user.UserId, "\tuser_pass:", user.Password)
	fmt.Println("User number:", user.UserNumber)

    w.Header().Set("HX-Location", "/dashboard") 
    w.WriteHeader(http.StatusSeeOther)
	fmt.Println("\tUser verify")
}

func RenderTemplate(w http.ResponseWriter, path string, data map[string]string) {
	tmpl, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/"+path+".html",
	)
	if err != nil {
		fmt.Printf("err.Error: %v\n", err.Error())
		http.Error(w, "Error loading template on "+path, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("err.Error: %v\n", err.Error())
		http.Error(w, "Error loading template on "+path, http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	content := map[string]string{
		"title":     "Login - ABCo-op",
		"stylePath": "login.css",
	}

	RenderTemplate(w, "login", content)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    content := map[string]string{
		"title":     "Register - ABCo-op",
		"stylePath": "register.css",
        "currentDate": time.Now().Format("2006-01-02"),
        "minDate": time.Now().AddDate(-100, 0, 0).Format("2006-01-02"),
    }

	RenderTemplate(w, "register", content)
}

func main() {
	static := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", static))

	mux.HandleFunc("/register-user/", registerUser)
	mux.HandleFunc("/verify-user/", verifyUser)
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/login", LoginHandler)

	fmt.Println("Server started at:\nlocalhost:5312/")
	log.Fatal(http.ListenAndServe(":5412", mux))
}

