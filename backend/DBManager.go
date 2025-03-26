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

func GetPaymentIdOf(loanId string, paymentNumber string) (int, error) {
	con, err := getConnection()
	if err != nil {
		return 0, err
	}
	defer con.Close()
	var paymentId int
	query := `SELECT PAYMENT_ID FROM PAYMENTS WHERE PAYMENTS.LOAN_ID = ? AND PAYMENTS.PAYMENT_NUMBER = ?`
	err = con.Get(&paymentId, query, loanId, paymentNumber)
	if err != nil {
		return 0, fmt.Errorf("Crash while fetching payment id!\nerr.Error(): %v\n", err.Error())
	}
	return paymentId, nil
}

func FetchClosures() ([]Closures, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	query := `SELECT *, closure_id || ': '||closure_month || '/' || closure_year AS closure_compact FROM closures;`
	closure := []Closures{}
	err = con.Select(&closure, query)
	if err != nil {
		return nil, fmt.Errorf("Crash while fetching the closures!\nerr.Error(): %v\n", err.Error())
	}

	return closure, nil
}

func FetchPayouts(closureId int) ([]Payouts, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	query := `
        SELECT 
            p.*, 
            APPORTATION_PERCENTAGE/100.00 AS decimal_percentage,
            u.FIRST_NAME || ' ' || u.FIRST_LASTNAME AS Name
        FROM payouts p 
        JOIN ACCOUNTS a 
        ON a.ACCOUNT_ID = p.ACCOUNT_ID
        JOIN USERS u
        ON a.USER_ID = u.USER_ID 
        WHERE p.closure_id = ?;`
	payouts := []Payouts{}
	err = con.Select(&payouts, query, closureId)
    if err != nil {
		return nil, fmt.Errorf("Crash while fetching the dividends!\nerr.Error(): %v\n", err.Error())
	}
    
    totalPayout := Payouts{
        Name: "TOTAL",
    }
    for i := range payouts {
        totalPayout.AccountBalance += payouts[i].AccountBalance
        totalPayout.AccountProfit += payouts[i].AccountProfit
        totalPayout.DecimalPercentage += payouts[i].DecimalPercentage
    }
    payouts = append(payouts, totalPayout)

    return payouts, nil
}

func FetchPayments(loanId string) ([]Payments, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	query := `SELECT 
            p.PAYMENT_ID, 
                MAX(p.LOAN_ID) LOAN_ID, 
                MAX(p.PAYMENT_NUMBER) PAYMENT_NUMBER, 
                MAX(p.DEADLINE) DEADLINE, 
                MAX(p.INTEREST_TO_PAY) INTEREST_TO_PAY, 
                MAX(p.CAPITAL_TO_PAY) CAPITAL_TO_PAY, 
                MAX(p.AMOUNT_TO_PAY) AMOUNT_TO_PAY, 
                MAX(p.IS_PAYED) IS_PAYED, 
                COALESCE(SUM(t.TRANSACTION_AMMOUNT), 0) AS AMOUNT_PAYED,
                (MAX(p.AMOUNT_TO_PAY) - COALESCE(SUM(t.TRANSACTION_AMMOUNT), 0)) AS REMAINING_AMOUNT
            FROM payments p
                LEFT JOIN payment_transactions pt
                    ON p.PAYMENT_ID = pt.PAYMENT_ID
                LEFT JOIN transactions t
                    ON pt.TRANSACTION_ID = t.TRANSACTION_ID
            WHERE p.LOAN_ID = ? AND p.IS_PAYED = FALSE
                GROUP BY p.PAYMENT_ID`
	payments := []Payments{}
	err = con.Select(&payments, query, loanId)
	if err != nil {
		return nil, fmt.Errorf("Error while getting the list of payments!\nerr.Error(): %v\n", err.Error())
	}

	if len(payments) == 0 {
		return nil, fmt.Errorf("There are no payments left!\nTo access this you'll need a loan first!\n")
	}

	for i := range payments {
		payments[i].FmtDeadline = payments[i].Deadline.Format("2006-02-03")
	}

	return payments, nil
}

func GetBalanceOf(id string) (float64, error) {
	con, err := getConnection()
	if err != nil {
		return 10000, err
	}
	defer con.Close()
	query := `SELECT balance FROM accounts WHERE account_id = ?`
	var balance float64
	err = con.Get(&balance, query, id)
	if err != nil {
		return 10000, fmt.Errorf("Crash while fetching the balance on account!\nerr.Error(): %v\n", err.Error())
	}
	return balance, nil
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

func GetLoanIdOfUser(id string) (string, error) {
	conn, err := getConnection()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var loanId string
	err = conn.Get(&loanId, `SELECT loan_id FROM loans l WHERE l.USER_ID = ? AND l.IS_PAYED = FALSE`, id)
	if err != nil {
		return "", fmt.Errorf("Couldn't fetch the loan id, please check if you have one active on the 'loans' tab!  error(): %v\n", err.Error())
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

func Update(crud Crudeable) error {
	con, err := getConnection()

	if err != nil {
		return err
	}
	defer con.Close()

	return crud.Update(con)
}

func FetchClosureById(Closure *Closures) error {
    con, err := getConnection()
	if err != nil {
		return err
	}
	defer con.Close()

    query := `SELECT * FROM closures WHERE closure_id = ?`;
    err = con.Get(Closure, query, Closure.Id)
    if err != nil {
        return fmt.Errorf("Crash while fetching closure by id!\nerr.Error(): %v\n", err.Error())
    }
    return nil
}

