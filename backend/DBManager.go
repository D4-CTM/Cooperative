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

func GrantAdminTo(id string, admin bool) error {
	con, err := getConnection()
	defer con.Close()
	if err != nil {
		return err
	}
	query := `CALL sp_grant_admin_to(?, ?)`

	_, err = con.Exec(query, id, admin)
	if err != nil {
		return fmt.Errorf("There was an error granting admin to %s\nError(): %s\n", id, err.Error())
	}

	return nil
}

func GetLoanIdOfUser(id string) (string, error) {
    conn, err := getConnection()
    if err !=nil {
        return "", err
    }

    var loanId string
	err = conn.Get(&loanId, `SELECT loan_id FROM loans l WHERE l.USER_ID = ? AND l.IS_PAYED = FALSE`, id)
	if err != nil {
		return "", fmt.Errorf("Crash while getting the loan_id\nerr.Error(): %v\n", err.Error())
	}

	return loanId, nil
}

func Fetch(fetchable Fetchable) error {
	con, err := getConnection()
	defer con.Close()

	if err != nil {
		return err
	}

	return fetchable.Fetch(con)
}

func Insert(crud Crudeable) error {
	con, err := getConnection()
	defer con.Close()

	if err != nil {
		return err
	}

	return crud.Insert(con)
}
