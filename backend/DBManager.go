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

func FetchLoanTransactionsYear(userId string) ([]int, error) {
    con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
    query := `SELECT 
            EXTRACT(YEAR FROM t.TRANSACTION_DATE)
        FROM LOANS l
        JOIN PAYMENTS p
        ON l.LOAN_ID = p.LOAN_ID
        JOIN PAYMENT_TRANSACTIONS pt 
        ON pt.PAYMENT_ID = p.PAYMENT_ID
        JOIN TRANSACTIONS t 
        ON t.TRANSACTION_ID = pt.TRANSACTION_ID
        WHERE l.USER_ID = ?
        GROUP BY EXTRACT(YEAR FROM t.TRANSACTION_DATE)
        ORDER BY EXTRACT(YEAR FROM t.TRANSACTION_DATE) DESC`
    year := []int{}
    err = con.Select(&year, query, userId)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching the years of loan transactions!\nerr.Error() %v\n", err.Error())
    }

    if len(year) == 0 {
        return nil, fmt.Errorf("Crash while fetching the years of loan transactions!\nerr.Error(): did'n found any transaction")
    }

    return year, nil
}

func FetchLoanTransactions(userId string, year int) ([]LoanTransactions, error) {
    con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
    query := `SELECT 
            l.LOAN_ID,
            pt.TRANSACTION_ID,
            t.TRANSACTION_AMMOUNT,
            t.TRANSACTION_DATE,
            p.PAYMENT_NUMBER
        FROM LOANS l
        JOIN PAYMENTS p
        ON l.LOAN_ID = p.LOAN_ID
        JOIN PAYMENT_TRANSACTIONS pt 
        ON pt.PAYMENT_ID = p.PAYMENT_ID
        JOIN TRANSACTIONS t 
        ON t.TRANSACTION_ID = pt.TRANSACTION_ID
        WHERE l.USER_ID = ? AND
        EXTRACT(YEAR FROM l.LOAN_DATE) = ?`

    loanTransactions := []LoanTransactions{}
    err = con.Select(&loanTransactions, query, userId, year)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching transactions from %d of the account: %s!\nerr.Error() %v\n", year, userId, err.Error())
    }

    if len(loanTransactions) == 0 {
        return nil, fmt.Errorf("Crash while fetching the yearly transactions of account: %s!\nerr.Error(): did'n found any transactions", userId)
    }

    totalLoanTransaction := LoanTransactions{
        LoanId: "TOTAL",
        TransactionId: fmt.Sprint(len(loanTransactions)),
    }
    for i := range loanTransactions {
        totalLoanTransaction.Amount += loanTransactions[i].Amount
        loanTransactions[i].FmtDate = loanTransactions[i].Date.Format("2006-01-02")
    }
    loanTransactions = append(loanTransactions, totalLoanTransaction)

    return loanTransactions, nil
}

func FetchTransactionsByYear(accId string, year int) ([]Transactions, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
    query := `SELECT *
        FROM TRANSACTIONS t 
        WHERE t.ACCOUNT_ID = ?
        AND EXTRACT(YEAR FROM t.TRANSACTION_DATE) = ?
        ORDER BY EXTRACT(YEAR FROM t.TRANSACTION_DATE) DESC;`
    transactions := []Transactions{}
    err = con.Select(&transactions, query, accId, year)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching transactions from %d of the account: %s!\nerr.Error() %v\n", year, accId, err.Error())
    }

    if len(transactions) == 0 {
        return nil, fmt.Errorf("Crash while fetching the yearly transactions of account: %s!\nerr.Error(): did'n found any transactions", accId)
    }

    totalTransaction := Transactions{
        TransactionId: "TOTAL",
        Amount: 0,
    }
    for i := range transactions {
        totalTransaction.Amount += transactions[i].Amount
        transactions[i].FmtDate = transactions[i].Date.Format("2006-01-02")
    }
    transactions = append(transactions, totalTransaction)

    return transactions, nil
}

func FetchTransactionsYears(accType string, accId string) ([]int, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
    query := `SELECT EXTRACT(YEAR FROM t.TRANSACTION_DATE)
        FROM TRANSACTIONS t
        JOIN ACCOUNTS a
        ON t.ACCOUNT_ID = a.ACCOUNT_ID
        WHERE a.USER_ID = ? AND a.ACCOUNT_TYPE = ?
        GROUP BY EXTRACT(YEAR FROM t.TRANSACTION_DATE)
        ORDER BY EXTRACT(YEAR FROM t.TRANSACTION_DATE) DESC`
    year := []int{}
    err = con.Select(&year, query, accId, accType)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching the years of transactions!\nerr.Error() %v\n", err.Error())
    }

    if len(year) == 0 {
        return nil, fmt.Errorf("Crash while fetching the years of transactions!\nerr.Error(): did'n found any transaction")
    }

    return year, nil

}

func FetchNewAccountYears() ([]int, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
    query := `SELECT EXTRACT(YEAR FROM HIRING_DATE)
        FROM USERS
        GROUP BY EXTRACT(YEAR FROM HIRING_DATE)
        ORDER BY EXTRACT(YEAR FROM HIRING_DATE) DESC`
    year := []int{}
    err = con.Select(&year, query)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching the affiliate report years!\nerr.Error() %v\n", err.Error())
    }

    if len(year) == 0 {
        return nil, fmt.Errorf("Crash while fetching the afiliate report years!\nerr.Error(): did'n found any affiliats")
    }

    return year, nil
}

func FetchAccountsReportInYear(year int) ([]AffiliateReports, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()

    query := `SELECT 
            MAX(u.USER_ID) AS user_id,
            MAX(u.FIRST_NAME) || ' ' || MAX(u.FIRST_LASTNAME) AS name,
            MAX(u.HIRING_DATE) AS hiring_date,
            MAX(CASE WHEN a.ACCOUNT_TYPE = 'CAR' THEN a.BALANCE END) AS apportation_balance,
            MAX(CASE WHEN a.ACCOUNT_TYPE = 'CAP' THEN a.BALANCE END) AS savings_balance,
            SUM(a.BALANCE) AS total
        FROM USERS u
        JOIN ACCOUNTS a
        ON u.USER_ID = a.USER_ID
        WHERE EXTRACT(YEAR FROM HIRING_DATE) = ?
        GROUP BY u.USER_ID`
    affiliateReports := []AffiliateReports{}
    err = con.Select(&affiliateReports, query, year)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching yearly affiliete report\nerr.Error(): %v\n", err.Error())
    }

    if len(affiliateReports) == 0 {
        return nil, fmt.Errorf("Crash while fetching yearly affiliete report\nerr.Error(): Didn't find any new affiliats during the year %v\n", year)
    }

    totalAffiliates := AffiliateReports{
        Name: "TOTAL",
        UserId: fmt.Sprint(len(affiliateReports)),
        HiringDateFmt: "",
    }

    for i := range affiliateReports {
        totalAffiliates.ApportationBalance += affiliateReports[i].ApportationBalance
        totalAffiliates.SavingsBalance += affiliateReports[i].SavingsBalance
        totalAffiliates.Total += affiliateReports[i].Total
        affiliateReports[i].HiringDateFmt = affiliateReports[i].HiringDate.Format("2006-01-02")
    }

    affiliateReports = append(affiliateReports, totalAffiliates)
    return affiliateReports, nil
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

func FetchAccountPayoutsYears(accountId string) ([]int, error){
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
    query := `SELECT c.CLOSURE_YEAR
        FROM PAYOUTS p 
        JOIN CLOSURES c
        ON p.CLOSURE_ID = c.CLOSURE_ID
        WHERE p.ACCOUNT_ID = ?
        GROUP BY CLOSURE_YEAR 
        ORDER BY c.CLOSURE_YEAR DESC;`
    year := []int{}
    err = con.Select(&year, query, accountId)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching the payout years!\nerr.Error() %v\n", err.Error())
    }

    if len(year) == 0 {
        return nil, fmt.Errorf("Crash while fetching the payout years!\nerr.Error(): did'n found any payout")
    }

    return year, nil
}

func FetchAccountPayouts(accountId string, year int) ([]Payouts, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	query := `
        SELECT 
            p.*,
            p.APPORTATION_PERCENTAGE/100.00 AS DECIMAL_PERCENTAGE,
            c.CLOSURE_YEAR || '/' || c.CLOSURE_MONTH AS PAYOUT_DATE
        FROM PAYOUTS p 
        JOIN CLOSURES c 
        ON c.CLOSURE_ID = p.CLOSURE_ID
        WHERE p.ACCOUNT_ID = ?
        AND c.CLOSURE_YEAR = ?`
	payouts := []Payouts{}
	err = con.Select(&payouts, query, accountId, year)
    if err != nil {
        return nil, fmt.Errorf("Crash while fetching the yearly payouts of user: %s!\nerr.Error(): %v\n", accountId, err.Error())
	}
   
    if len(payouts) == 0 {
        return nil, fmt.Errorf("Crash while fetching the user specific yearly payouts!\nerr.Error(): Couldn't find any payout for this year!\n")
    }

    totalPayout := Payouts{
        PayoutDate: "TOTAL",
    }
    for i := range payouts {
        totalPayout.AccountBalance += payouts[i].AccountBalance
        totalPayout.AccountProfit += payouts[i].AccountProfit
        totalPayout.DecimalPercentage += payouts[i].DecimalPercentage
    }
    payouts = append(payouts, totalPayout)

    return payouts, nil
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
		payments[i].FmtDeadline = payments[i].Deadline.Format("2006-01-02")
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

