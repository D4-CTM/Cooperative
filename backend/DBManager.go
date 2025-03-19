package backend

import (
	"fmt"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/jmoiron/sqlx"
)

type LoginData struct {
	UserId   string
	Password string
	Name     string
	Admin    bool
	LoanId   string
}

var LoginUser LoginData = LoginData{UserId: ""}

func getConnection() (*sqlx.DB, error) {
	const CONNECTION_STRING string = "HOSTNAME=localhost;DATABASE=coopdb;PORT=51000;UID=db2inst1;PWD=coop4312"
	db, err := sqlx.Connect("go_ibm_db", CONNECTION_STRING)
	if err != nil {
		return nil, fmt.Errorf("Crashed on connection!\nerr.Error():%v\n", err.Error())
	}
	return db, nil
}

func FetchPayments(loanId string) ([]Payments, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
    defer con.Close()
    query := `SELECT * FROM payments WHERE loan_id = ?`
    payments := []Payments{}
    err = con.Select(&payments, query, loanId)
    if err != nil {
        return nil, fmt.Errorf("Error while getting the list of payments!\nerr.Error(): %v\n", err.Error()) 
    }

    for i := 0; i < len(payments); i++ {
        payments[i].FmtDeadline = payments[i].Deadline.Format("2006-02-03")
    }

    return payments, nil
}

func GrantAdminTo(id string, admin bool) error {
	con, err := getConnection()
	if err != nil {
		return err
	}
    defer con.Close()
	query := `CALL sp_grant_admin_to(?, ?)`

	_, err = con.Exec(query, id, admin)
	if err != nil {
		return fmt.Errorf("There was an error granting admin to %s\nError(): %s\n", id, err.Error())
	}

	return nil
}

func MarkLoanAsPayed(id string) error {
	con, err := getConnection()
	if err != nil {
		return err
	}
    defer con.Close()
	query := `CALL sp_mark_loan_as_payed(?)`

	_, err = con.Exec(query, id)
	if err != nil {
		return fmt.Errorf("There was an error marking the loan to %s\nError(): %s\n", id, err.Error())
	}

	return nil
}

func GetLoanIdOfUser(id string) (string, error) {
    conn, err := getConnection()
    if err !=nil {
        return "", err
    }
    defer conn.Close()

    var loanId string
	err = conn.Get(&loanId, `SELECT loan_id FROM loans l WHERE l.USER_ID = ? AND l.IS_PAYED = FALSE`, id)
	if err != nil {
		return "", fmt.Errorf("Crash while getting the loan_id\nerr.Error(): %v\n", err.Error())
	}

	return loanId, nil
}

func Fetch(fetchable Fetchable) error {
	con, err := getConnection()
	if err != nil {
		return err
	}
    defer con.Close()

	return fetchable.Fetch(con)
}

func Insert(crud Crudeable) error {
	con, err := getConnection()

	if err != nil {
		return err
	}
    defer con.Close()

	return crud.Insert(con)
}
