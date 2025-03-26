package main

import (
	"cooperative/backend"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func revClosure(w http.ResponseWriter, r *http.Request) {    
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

    if r.Header.Get("HX-Request") == "" {
        fmt.Println("\tWASN'T A HX-Request!")
    }
    
    closureId, err := strconv.Atoi(r.PostFormValue("closure-select"))
	
	if err != nil {
        fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return ;
    }

    if closureId < 1 {
		w.Header().Set("HX-Status", "202")
		w.WriteHeader(http.StatusAccepted)
        return ;
    }

    ct := backend.ClosureTransaction{
        ClosureId: closureId,
    }
    err = backend.Fetch(&ct)    
	if err != nil {
        fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return ;
    }

    cp := backend.ClosurePayments{
        ClosureId: closureId,
    }
    err = backend.Fetch(&cp)    
	if err != nil {
        fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return ;
    }
    
    payouts, err := backend.FetchPayouts(closureId)
	if err != nil {
        fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return ;
    }

    content := map[string]any {
        "ClosureTransactions": ct.Transactions,
        "ClosurePayments": cp.CPT,
        "Dividends": payouts,
    }

    w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
    tmpl, err := template.ParseFiles("./templates/DashboardOptions/revClosure.html")
	if err != nil {
		fmt.Println(err.Error())
        return
	}
	tmpl.ExecuteTemplate(w, "review-data", content)
}

func createClosure(r *http.Request) backend.Closures {
    desc := r.PostFormValue("description")
    month_name, err := time.Parse("January", r.PostFormValue("month"))
    if err != nil {
        fmt.Println(err.Error())
    }
    
    year, err := strconv.Atoi(r.PostFormValue("year"))
    if err != nil {
        fmt.Println(err.Error())
    }
    
    closure := backend.Closures{
        Year: year,
        Month: int(month_name.Month()),
        Description: desc,
    }

    return closure
}

func modClosure(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	fmt.Println("Starting register closure...")
	
    if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
	}
    closure := createClosure(r)
    err := backend.Update(&closure)
	if err != nil {
        fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return ;
    }

    w.Header().Set("hx-status", "200")
	w.Header().Set("hx-message", "Comment changed succesfully! You can check on the review closure tab!")
	w.WriteHeader(http.StatusOK)
    fmt.Println("Closure made!")
}

func regClosure(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	fmt.Println("Starting register closure...")
	
    if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
	}

    closure := createClosure(r)
    err := backend.Insert(&closure)
	if err != nil {
        fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return ;
    }

    w.Header().Set("hx-status", "200")
	w.Header().Set("hx-message", "Monthly closure registered succesfully! You can check on the review closures tab to check!")
	w.WriteHeader(http.StatusOK)
    fmt.Println("Closure made!")
}

func requestPayment(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	fmt.Println("Starting payment request")

	if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
	}

	transaction := CreateTransaction(r)
	payNo := r.PostFormValue("payment-number")

	payId, err := backend.GetPaymentIdOf(r.PostFormValue("loan-id"), payNo)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
        return
    }

	payment := backend.Payments{
		PaymentId: payId,
		LoanId:    r.PostFormValue("loan-id"),
	}

	PT := backend.PaymentTransaction{
		Payment:         payment,
		TransactionList: []backend.Transactions{transaction},
	}
	fmt.Println(PT)
	err = backend.Insert(&PT)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("hx-status", "200")
	w.Header().Set("hx-message", "Payment done! check your loan data to confirm it")
	w.WriteHeader(http.StatusOK)
}

func requestDeposit(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	fmt.Println("Starting deposit request")

	if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
	}

	transaction := CreateTransaction(r)

	err := backend.Insert(&transaction)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	    return
    }

	w.Header().Set("hx-status", "200")
	w.Header().Set("hx-message", "Deposit completed, check your account balance!")
	w.WriteHeader(http.StatusOK)
}

func CreateTransaction(r *http.Request) backend.Transactions {
	amount, err := strconv.ParseFloat(strings.TrimSpace(r.PostFormValue("amount")), 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	description := r.PostFormValue("description")

	transaction := backend.Transactions{
		Amount:    amount,
		Comment:   sql.NullString{String: description, Valid: len(description) > 0},
		AccountId: backend.LoginUser.UserId + r.PostFormValue("destiny-acc"),
		Date:      time.Now(),
	}

	return transaction
}

func requestLoan(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	fmt.Println("Starting loan request...")

	if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
	}

	capital, err := strconv.ParseFloat(strings.TrimSpace(r.PostFormValue("capital")), 64)
	if err != nil {
		fmt.Println("Crash while parsing capital\nerr.Error():", err.Error())
		DashboardLoanHandler(w, r)
		return
	}

	var interest float32 = 0.15
	if r.PostFormValue("loan-type") == "automatic" {
		interest = 0.10
	}

	periods, err := strconv.Atoi(r.PostFormValue("periods"))
	if err != nil {
		fmt.Println("Crash while parsing periods\nerr.Error():", err.Error())
		return
	}

	var loan backend.Loans = backend.Loans{
		UserId:   backend.LoginUser.UserId,
		Capital:  capital,
		Periods:  periods,
		Interest: interest,
	}

	err = backend.Insert(&loan)
	if err != nil {
		fmt.Println(err.Error())
		DashboardLoanHandler(w, r)
		return
	}

	fmt.Println("Loan id:", loan.LoanId)
	DashboardLoanHandler(w, r)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	fmt.Println("Starting registration user...")

	if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
		return
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

	user := backend.Users{
		Password:             _password,
		FirstName:            first_name,
		FirstLastname:        first_last_name,
		SecondName:           sql.NullString{String: second_name, Valid: len(strings.TrimSpace(second_name)) > 0},
		SecondLastname:       sql.NullString{String: second_lastname, Valid: len(strings.TrimSpace(second_lastname)) > 0},
		BirthDate:            sql.NullTime{Time: birthdate, Valid: birthdate.Year() > 1},
		PrimaryEmail:         email,
		SecondaryEmail:       sql.NullString{String: recovery_email, Valid: len(strings.TrimSpace(recovery_email)) > 0},
		AddressDepartment:    sql.NullString{String: department, Valid: len(strings.TrimSpace(department)) > 0},
		AddressCity:          sql.NullString{String: city, Valid: len(strings.TrimSpace(city)) > 0},
		AddressStreet:        sql.NullString{String: street, Valid: len(strings.TrimSpace(street)) > 0},
		AddressAvenue:        sql.NullString{String: avenue, Valid: len(strings.TrimSpace(avenue)) > 0},
		AddressHouseNumber:   sql.NullString{String: house_number, Valid: len(strings.TrimSpace(house_number)) > 0},
		AddressReference:     sql.NullString{String: reference, Valid: len(strings.TrimSpace(reference)) > 0},
		HiringDate:           time.Now(),
		CreatedBy:            sql.NullString{String: "Admin", Valid: true},
		CreationDate:         time.Now(),
		ModifiedBy:           sql.NullString{String: "Admin", Valid: true},
		LastModificationDate: time.Now(),
	}

	err := backend.Insert(&user)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "There was an error while creating your accound, please check that fill all the fields with an '*'")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("your id is:", user.UserId)

	w.Header().Set("HX-Location", "/login")
	w.Header().Set("HX-Status", "200")
	w.Header().Set("HX-Message", fmt.Sprintf("Affiliated succesfully! Your id is: %s", user.UserId))
	w.WriteHeader(http.StatusOK)
	fmt.Println("\tUser registered")
}

func verifyUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\tVerifying user...")
	if r.Header.Get("HX-Request") == "" {
		fmt.Println("\tWASN'T A HX-Request!")
		return
	}

	id := r.PostFormValue("code")
	pass := r.PostFormValue("password")

	user := backend.Users{
		UserId:   id,
		Password: pass,
	}

	err := backend.Fetch(&user)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<p>There was an error fetching the user!<br>Please check the user or password</p>`))
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("user_id:", user.UserId, "\tuser_pass:", user.Password)
	backend.LoginUser.UserId = user.UserId
	backend.LoginUser.Password = user.Password
	backend.LoginUser.Name = user.FirstName + " " + user.FirstLastname
	backend.LoginUser.Admin = user.Admin

	w.Header().Set("HX-Location", "/dashboard")
	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusOK)
	fmt.Println("\tUser verify")
}

func RenderTemplate(w http.ResponseWriter, data map[string]any, extraFiles ...string) {
	files := []string{
		"./templates/layout.html",
	}

	files = append(files, extraFiles...)

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		fmt.Printf("err.Error: %v\n", err.Error())
		http.Error(w, "Error loading templates:", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("err.Error: %v\n", err.Error())
		http.Error(w, "Error loading templates:", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	backend.LoginUser.UserId = ""

	content := map[string]any{
		"title":     "Login - ABCo-op",
		"stylePath": "login.css",
	}

	RenderTemplate(w, content, "./templates/login.html")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	backend.LoginUser.UserId = ""

	content := map[string]any{
		"title":       "Register - ABCo-op",
		"stylePath":   "register.css",
		"currentDate": time.Now().Format("2006-01-02"),
		"minDate":     time.Now().AddDate(-100, 0, 0).Format("2006-01-02"),
	}

	RenderTemplate(w, content, "./templates/register.html")
}

func stringBool(admin bool) string {
	if admin {
		return "T"
	}

	return "F"
}

func compareDates(t1 time.Time, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func DashboardUserHandler(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		fmt.Println("Please login first")
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	now := time.Now()
	year, month, _ := now.Date()
	lastOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1)

	content := map[string]any{
		"title":        "Dashboard - ABCo-op",
		"stylePath":    "dashboard.css",
		"user":         backend.LoginUser.Name,
		"admin":        stringBool(backend.LoginUser.Admin),
		"closureValid": stringBool(compareDates(now, lastOfMonth) || true),
	}

	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
	RenderTemplate(w, content, "./templates/dashboard.html", "./templates/DashboardOptions/user.html")
}

func DashboardLoanHandler(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		fmt.Println("Please login first")
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	loanId, err := backend.GetLoanIdOfUser(backend.LoginUser.UserId)
	if err != nil {
		fmt.Println(err.Error())
	}
	content := map[string]any{
		"admin":      stringBool(backend.LoginUser.Admin),
		"loanActive": stringBool(loanId != ""),
	}

	if content["loanActive"] == "T" {
		payments, err := backend.FetchPayments(loanId)
		if err != nil {
			fmt.Println(err.Error())
		}
		content["payments"] = payments

		loan := backend.Loans{LoanId: loanId}
		err = backend.Fetch(&loan)
		if err != nil {
			fmt.Println(err.Error())
		}
		// We change the UserId to the name of the user to make it more user friendly
		loan.UserId = backend.LoginUser.Name
		content["loanDate"] = loan.Date.Format("2006-02-03")
		content["loan"] = loan

	}
	maxAmount, err := backend.GetBalanceOf(backend.LoginUser.UserId + "-CAP")
	if err != nil {
		fmt.Println(err.Error())
		maxAmount = 0
	}
	content["MaxAmount"] = maxAmount

	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)

	tmpl, err := template.ParseFiles("./templates/DashboardOptions/loan.html")
	if err != nil {
		fmt.Println(err.Error())
	    return
    }
	tmpl.ExecuteTemplate(w, "dashboard-content", content)
}

func DashboardDepositHandle(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		fmt.Println("Please login first")
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	content := map[string]any{
		"admin":     stringBool(backend.LoginUser.Admin),
		"type":      "Deposit",
		"endpoint":  "/request-deposit/",
		"action":    "deposit",
		"MinAmount": 200,
	}
	content["MaxAmount"] = 10000
	content["UserName"] = backend.LoginUser.Name
	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
	tmpl, err := template.ParseFiles("./templates/DashboardOptions/transaction.html")
	if err != nil {
		fmt.Println(err.Error())
	    return
    }
	tmpl.ExecuteTemplate(w, "dashboard-content", content)
}

func DashboardPaymentHandle(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		fmt.Println("Please login first")
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	content := map[string]any{
		"admin":     stringBool(backend.LoginUser.Admin),
		"type":      "Payment",
		"endpoint":  "/request-payment/",
		"action":    "pay",
		"MinAmount": 0.01,
	}
	content["UserName"] = backend.LoginUser.Name
	loanId, err := backend.GetLoanIdOfUser(backend.LoginUser.UserId)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	payments, err := backend.FetchPayments(loanId)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	content["MaxAmount"] = payments[0].AmountToPay
	content["payments"] = payments
	content["loanId"] = loanId

	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
	tmpl, err := template.ParseFiles("./templates/DashboardOptions/transaction.html")
	if err != nil {
		fmt.Println(err.Error())
	    return
    }
	tmpl.ExecuteTemplate(w, "dashboard-content", content)
}

func DashboardLiquidationHandle(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	content := map[string]any{
		"admin":     stringBool(backend.LoginUser.Admin),
		"type":      "Liquidation",
		"endpoint":  "/request-liquidation/",
		"action":    "retire",
		"MinAmount": 0.01,
	}
	balance, err := backend.GetBalanceOf(backend.LoginUser.UserId + "-CAR")
	if err != nil {
		fmt.Println(err.Error())
	}
	content["MaxAmount"] = balance
	content["UserName"] = backend.LoginUser.Name

	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
	tmpl, err := template.ParseFiles("./templates/DashboardOptions/transaction.html")
	if err != nil {
		fmt.Println(err.Error())
	    return
    }
	tmpl.ExecuteTemplate(w, "dashboard-content", content)
}

func DashboardRegisterClosure(w http.ResponseWriter, r *http.Request) {
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

    endPoint := "/modify-closure/"
    action := "Modify"
    year, month, _ := time.Now().Date()
    closure := backend.Closures{
        Year: year,
        Month: int(month),
    }
    err := backend.Fetch(&closure)
    if err != nil {
        fmt.Println(err.Error())
        endPoint = "/register-closure/"
        action = "Register"
    }

	content := map[string]any{
		"year":  year,
		"month": month,
        "endPoint": endPoint,
        "action": action,
    }

	w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
	tmpl, err := template.ParseFiles("./templates/DashboardOptions/regClosure.html")
	if err != nil {
		fmt.Println(err.Error())
        return
	}
	tmpl.ExecuteTemplate(w, "dashboard-content", content)
}

func DashboardReviewClosures(w http.ResponseWriter, r *http.Request) {    
	if len(backend.LoginUser.UserId) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Please login first")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
    
    closures, err := backend.FetchClosures()
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

    if len(closures) == 0 {
		w.Header().Set("HX-Status", "400")
		w.Header().Set("HX-Message", "Couldn't find any closure registered!")
		w.WriteHeader(http.StatusNotFound)
		return
    }

    content := map[string]any {
        "Closures": closures,
    }

    w.Header().Set("HX-Status", "202")
	w.WriteHeader(http.StatusAccepted)
    tmpl, err := template.ParseFiles("./templates/DashboardOptions/revClosure.html")
	if err != nil {
		fmt.Println(err.Error())
        return
	}
	tmpl.ExecuteTemplate(w, "dashboard-content", content)
}

func main() {
	static := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", static))

    // basically requests
    mux.HandleFunc("/review-closure/", revClosure)
    mux.HandleFunc("/modify-closure/", modClosure)
	mux.HandleFunc("/register-closure/", regClosure)
	mux.HandleFunc("/request-payment/", requestPayment)
	mux.HandleFunc("/request-deposit/", requestDeposit)
	mux.HandleFunc("/request-loan/", requestLoan)
	mux.HandleFunc("/register-user/", registerUser)
	mux.HandleFunc("/verify-user/", verifyUser)
	// dashboard options
    mux.HandleFunc("/dashboard-rev-closure/", DashboardReviewClosures)
    mux.HandleFunc("/dashboard-reg-closure/", DashboardRegisterClosure)
    mux.HandleFunc("/dashboard-deposits/", DashboardDepositHandle)
	mux.HandleFunc("/dashboard-liquidations/", DashboardLiquidationHandle)
	mux.HandleFunc("/dashboard-payments/", DashboardPaymentHandle)
	mux.HandleFunc("/dashboard-user/", DashboardUserHandler)
	mux.HandleFunc("/dashboard-loan/", DashboardLoanHandler)
    // general gui options
    mux.HandleFunc("/dashboard", DashboardUserHandler)
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/login", LoginHandler)

	fmt.Println("Server started at:\nlocalhost:5312/")
	log.Fatal(http.ListenAndServe(":5412", mux))
}

