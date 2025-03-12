package main

import (
	"cooperative/backend"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func registerUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Starting registration user...")

    if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
        return ;
    }

    first_name := r.PostFormValue("first_name")
    second_name := r.PostFormValue("second_name")
    first_last_name := r.PostFormValue("first_lastname")
    second_lastname := r.PostFormValue("second_lastname")
    birthdate, _ := time.Parse("2006-01-02", r.PostFormValue("birth_date")) 
    email := r.PostFormValue("email")
    recovery_email := r.PostFormValue("recovery_email")
    _password := r.PostFormValue("password")
    department := r.PostFormValue("department")
    city := r.PostFormValue("city")
    street := r.PostFormValue("street")
    avenue := r.PostFormValue("avenue")
    house_number := r.PostFormValue("house_number")
    reference := r.PostFormValue("reference")

    user := backend.User {
        Password: _password,
        FirstName: first_name,
        FirstLastname: first_last_name,
        SecondName: sql.NullString{String: second_name, Valid: len(strings.TrimSpace(second_name)) > 0},
        SecondLastname: sql.NullString{String: second_lastname, Valid: len(strings.TrimSpace(second_lastname)) > 0},
        BirthDate: sql.NullTime{Time: birthdate, Valid: birthdate.Year() > 1},
        PrimaryEmail: email,
        SecondaryEmail: sql.NullString{String: recovery_email, Valid: len(strings.TrimSpace(recovery_email)) > 0},
        AddressDepartment: sql.NullString{String: department, Valid: len(strings.TrimSpace(department)) > 0},
        AddressCity: sql.NullString{String: city, Valid: len(strings.TrimSpace(city)) > 0},
        AddressStreet: sql.NullString{String: street, Valid: len(strings.TrimSpace(street)) > 0},
        AddressAvenue: sql.NullString{String: avenue, Valid: len(strings.TrimSpace(avenue)) > 0},
        AddressHouseNumber: sql.NullString{String: house_number, Valid: len(strings.TrimSpace(house_number)) > 0},
        AddressReference: sql.NullString{String: reference, Valid: len(strings.TrimSpace(reference)) > 0},
        HiringDate: time.Now(),
        CreatedBy: sql.NullString{String: "Admin", Valid: true},
        CreationDate: time.Now(),
        ModifiedBy: sql.NullString{String: "Admin", Valid: true},
        LastModificationDate: time.Now(),
    }

    err := backend.Insert(&user)
    if err != nil {
        fmt.Println(err.Error())
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
	fmt.Println("User id:", user.UserId)

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

