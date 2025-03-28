package backend

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Crudeable interface {
	Insert(db *sqlx.DB) error
	Update(db *sqlx.DB) error
}

type Fetchable interface {
	Fetch(db *sqlx.DB) error
}

type Users struct {
	UserId               string         `db:"USER_ID"`
	Password             string         `db:"PASSWORD"`
	Admin                bool           `db:"IS_ADMIN"`
	IsActive             bool           `db:"IS_ACTIVE"`
	FirstName            string         `db:"FIRST_NAME"`
	SecondName           sql.NullString `db:"SECOND_NAME"`
	FirstLastname        string         `db:"FIRST_LASTNAME"`
	SecondLastname       sql.NullString `db:"SECOND_LASTNAME"`
	AddressHouseNumber   sql.NullString `db:"ADDRESS_HOUSE_NUMBER"`
	AddressStreet        sql.NullString `db:"ADDRESS_STREET"`
	AddressAvenue        sql.NullString `db:"ADDRESS_AVENUE"`
	AddressCity          sql.NullString `db:"ADDRESS_CITY"`
	AddressDepartment    sql.NullString `db:"ADDRESS_DEPARTMENT"`
	AddressReference     sql.NullString `db:"ADDRESS_REFERENCE"`
	PrimaryEmail         string         `db:"PRIMARY_EMAIL"`
	SecondaryEmail       sql.NullString `db:"SECONDARY_EMAIL"`
	BirthDate            sql.NullTime   `db:"BIRTH_DATE"`
	HiringDate           time.Time      `db:"HIRING_DATE"`
	CreatedBy            sql.NullString `db:"CREATED_BY"`
	CreationDate         time.Time      `db:"CREATION_DATE"`
	ModifiedBy           sql.NullString `db:"MODIFIED_BY"`
	LastModificationDate time.Time      `db:"LAST_MODIFICATION_DATE"`
}

func (user *Users) Insert(db *sqlx.DB) error {
	query := `CALL sp_insert_user(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(query,
		sql.Out{Dest: &user.UserId},
		user.Password,
		user.FirstName,
		user.SecondName,
		user.FirstLastname,
		user.SecondLastname,
		user.AddressHouseNumber,
		user.AddressStreet,
		user.AddressAvenue,
		user.AddressCity,
		user.AddressDepartment,
		user.AddressReference,
		user.PrimaryEmail,
		user.SecondaryEmail,
		user.BirthDate,
		user.HiringDate,
		user.CreatedBy,
		user.CreationDate,
		user.ModifiedBy,
		user.LastModificationDate,
		user.Admin)

	if err != nil {
		return fmt.Errorf("Crash at user insert!\nerr.Error(): %v\n", err.Error())
	}

	fmt.Println("User inserted succesfully, id:", user.UserId)
	return nil
}

func (user *Users) Update(db *sqlx.DB) error {
	query := `CALL sp_update_user(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(query,
		sql.Out{Dest: &user.UserId},
		user.Password,
		user.FirstName,
		user.SecondName,
		user.FirstLastname,
		user.SecondLastname,
		user.AddressHouseNumber,
		user.AddressStreet,
		user.AddressAvenue,
		user.AddressCity,
		user.AddressDepartment,
		user.AddressReference,
		user.PrimaryEmail,
		user.SecondaryEmail,
		user.BirthDate,
		user.ModifiedBy,
		user.LastModificationDate,
		user.Admin)
	if err != nil {
		return fmt.Errorf("Crash while updating user\nerr.Error(): %v\n", err.Error())
	}
    
    fmt.Println("updated!")
	return nil
}

func (user *Users) Fetch(db *sqlx.DB) error {
	err := db.Get(user, `SELECT * FROM users WHERE user_id=? AND password=?`, user.UserId, user.Password)
	if err != nil {
		return fmt.Errorf("The user you are searching was not found!\nerr.Error(): %v\n", err.Error())
	}
	if !user.IsActive {
		return fmt.Errorf("The user you are searching was not found!\nerr.Error(): The user has retire from the company")
	}

	return nil
}

type PhoneNumbers struct {
	UserId          string `db:"USER_ID"`
	UserPhoneNumber int    `db:"USER_PHONE_NUMBER"`
	RegionNumber    int    `db:"REGION_NUMBER"`
}

func (pn *PhoneNumbers) Insert(db *sqlx.DB) error {
	query := `CALL sp_insert_phone(?, ?, ?)`

	_, err := db.Exec(query, pn.UserId, pn.UserPhoneNumber, pn.RegionNumber)
	if err != nil {
		return fmt.Errorf("Crash while inserting phone number!\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

func (pn *PhoneNumbers) Update(db *sqlx.DB) error {
	query := `CALL sp_update_phone(?, ?, ?)`

	_, err := db.Exec(query, pn.UserId, pn.UserPhoneNumber, pn.RegionNumber)
	if err != nil {
		return err
	}

	return nil
}

type Loans struct {
	LoanId               string    `db:"LOAN_ID"`
	UserId               string    `db:"USER_ID"`
	Periods              int       `db:"LOAN_PERIODS"`
	Interest             float32   `db:"LOAN_INTEREST"`
	Capital              float64   `db:"REQUESTED_AMOUNT"`
	Date                 time.Time `db:"LOAN_DATE"`
	IsPayed              bool      `db:"IS_PAYED"`
	CreatedBy            string    `db:"CREATED_BY"`
	CreationDate         time.Time `db:"CREATION_DATE"`
	Modified_by          string    `db:"MODIFIED_BY"`
	LastModificationDate time.Time `db:"LAST_MODIFICATION"`
}

func (loan *Loans) Insert(db *sqlx.DB) error {
	query := `CALL sp_create_loan(?,?,?,?,?,?,?)`
	_, err := db.Exec(query,
		sql.Out{Dest: &loan.LoanId},
		sql.Out{Dest: &loan.Date},
		sql.Out{Dest: &loan.IsPayed},
		loan.UserId,
		loan.Periods,
		loan.Interest,
		loan.Capital)

	if err != nil {
		return fmt.Errorf("Crash at loan insertion!\nerr.Error(): %v\n", err.Error())
	}

	fmt.Printf("Loan inserted succesfully!\nuser_id: %s\tloan_id: %s\n", loan.UserId, loan.LoanId)
	return nil
}

func (loan *Loans) Update(db *sqlx.DB) error {
	query := `CALL sp_update_loan(?,?,?,?,?)`
	_, err := db.Exec(query,
		loan.LoanId,
		loan.Periods,
		loan.Interest,
		loan.Capital,
		loan.Date)
	if err != nil {
		return fmt.Errorf("Crash while updating the loan!\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

func (loan *Loans) Fetch(db *sqlx.DB) error {
	err := db.Get(loan, `SELECT * FROM loans WHERE loan_id = ?`, loan.LoanId)
	if err != nil {
		return fmt.Errorf("Crash while fetching loans\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

type Payments struct {
	PaymentId     int       `db:"PAYMENT_ID"`
	LoanId        string    `db:"LOAN_ID"`
	PaymentNumber string    `db:"PAYMENT_NUMBER"`
	Deadline      time.Time `db:"DEADLINE"`
	IPMT          float64   `db:"INTEREST_TO_PAY"`
	PPMT          float64   `db:"CAPITAL_TO_PAY"`
	PMT           float64   `db:"AMOUNT_TO_PAY"`
	IsPayed       bool      `db:"IS_PAYED"`
	AmountPayed   float64   `db:"AMOUNT_PAYED"`
	AmountToPay   float64   `db:"REMAINING_AMOUNT"`
	FmtDeadline   string
}

func (payment *Payments) Fetch(db *sqlx.DB) error {
	err := db.Get(payment, `SELECT * FROM payments WHERE loan_id = ? AND payment_number = ?`, payment.LoanId, payment.PaymentNumber)
	if err != nil {
		return fmt.Errorf("Crash while fetching loans\nerr.Error(): %v\n", err.Error())
	}
	payment.FmtDeadline = payment.Deadline.Format("2006-02-03")
	return nil
}

type Transactions struct {
	AccountId     string         `db:"ACCOUNT_ID"`
	TransactionId string         `db:"TRANSACTION_ID"`
	Date          time.Time      `db:"TRANSACTION_DATE"`
	Amount        float64        `db:"TRANSACTION_AMMOUNT"`
	Comment       sql.NullString `db:"TRANSACTION_COMMENT"`
	FmtDate       string
}

func (transaction *Transactions) Insert(db *sqlx.DB) error {
	query := `CALL sp_create_transaction(?, ?, ?, ?, ?)`
	_, err := db.Exec(query,
		sql.Out{Dest: &transaction.TransactionId},
		transaction.AccountId,
		transaction.Date,
		transaction.Amount,
		transaction.Comment)
	if err != nil {
		return fmt.Errorf("Error when inserting transaction!\n %v\n", err.Error())
	}

	return nil
}

func (transaction Transactions) Update(db *sqlx.DB) error {
	query := `CALL sp_change_transaction_comment(?, ?)`
	_, err := db.Exec(query,
		transaction.TransactionId,
		transaction.Comment)
	if err != nil {
		return fmt.Errorf("Error when changing transaction comment!\n %v\n", err.Error())
	}

	return nil
}

func (transaction *Transactions) Fetch(db *sqlx.DB) error {
	err := db.Get(transaction, `SELECT * FROM transactions WHERE transaction_id = ?`, transaction.TransactionId)
	if err != nil {
		return fmt.Errorf("Crash while fetching transaction!\nerr.Error() %v\n", err.Error())
	}
	transaction.FmtDate = transaction.Date.Format("2006-02-03")
	return nil
}

type PaymentTransaction struct {
	Payment         Payments
	TransactionList []Transactions
}

func (PT *PaymentTransaction) Insert(db *sqlx.DB) error {
	query := `CALL sp_payment_transaction(?, ?, ?, ?, ?, ?, ?)`
	for i := range PT.TransactionList {
		fmt.Println("Transaction send:", PT.TransactionList[i])
		_, err := db.Exec(query,
			sql.Out{Dest: &PT.TransactionList[i].TransactionId},
			PT.TransactionList[i].AccountId,
			PT.TransactionList[i].Date,
			PT.TransactionList[i].Amount,
			PT.TransactionList[i].Comment,
			PT.Payment.PaymentId,
			PT.Payment.LoanId)

		if err != nil {
			return fmt.Errorf("Crash while inserting #%d payment transaction!\nerr.Error(): %v\n", i, err.Error())
		}

		fmt.Printf("\nTransaction #%d done!\n\tTransaction id: %s", i, PT.TransactionList[i].TransactionId)
	}

	return nil
}

func (PT *PaymentTransaction) Update(db *sqlx.DB) error {
	return fmt.Errorf("\nPayment transaction doesn't support updates!\n")
}

func (PT *PaymentTransaction) Fetch(db *sqlx.DB) error {
	query := `SELECT t.*
        FROM PAYMENT_TRANSACTIONS pt
        JOIN TRANSACTIONS t 
        ON pt.TRANSACTION_ID = t.TRANSACTION_ID
        WHERE pt.PAYMENT_ID = ?`
	err := db.Select(&PT.TransactionList, query, PT.Payment.PaymentId)
	if err != nil {
		return fmt.Errorf("Crash while fetching payment transactions!\nerr.Error(): %v\n", err.Error())
	}
	return nil
}

type Closures struct {
	Id          int    `db:"CLOSURE_ID"`
	Month       int    `db:"CLOSURE_MONTH"`
	Year        int    `db:"CLOSURE_YEAR"`
	Description string `db:"DESCRIPTION"`
	Compact     string `db:"CLOSURE_COMPACT"`
}

func (closure *Closures) Insert(db *sqlx.DB) error {
	query := `CALL sp_generate_monthly_closure(?, ?, ?)`
	_, err := db.Exec(query,
		closure.Month,
		closure.Year,
		closure.Description)

	if err != nil {
		return fmt.Errorf("Crash while generating the monthly closure!\nerr.Error(): %v\n", err.Error())
	}
	return nil
}

func (closure *Closures) Update(db *sqlx.DB) error {
	query := `CALL sp_change_closure_comment(?, ?, ?)`
	_, err := db.Exec(query,
		closure.Month,
		closure.Year,
		closure.Description)
	if err != nil {
		return fmt.Errorf("Crash while changing the monthly closure description!\nerr.Error(): %v\n", err.Error())
	}
	return nil
}

func (closure *Closures) Fetch(db *sqlx.DB) error {
	query := `SELECT *, closure_id || ': '||closure_month || '/' || closure_year AS closure_compact FROM CLOSURES WHERE closure_month = ? AND closure_year = ?`
	err := db.Get(closure,
		query,
		closure.Month,
		closure.Year)
	if err != nil {
		return fmt.Errorf("Crash while fetching the closure!\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

type ClosureTransaction struct {
	ClosureId    int
	Transactions []Transactions
}

func (closureTransaction *ClosureTransaction) Fetch(db *sqlx.DB) error {
	query := `SELECT t.* 
        FROM CLOSURE_TRANSACTIONS ct 
        JOIN TRANSACTIONS t 
        ON ct.TRANSACTION_ID = t.TRANSACTION_ID
        LEFT JOIN PAYMENT_TRANSACTIONS pt 
        ON pt.TRANSACTION_ID = t.TRANSACTION_ID
        WHERE ct.CLOSURE_ID = ?
        AND pt.PAYMENT_ID IS NULL;`

	err := db.Select(&closureTransaction.Transactions,
		query,
		closureTransaction.ClosureId)
	if err != nil {
		return fmt.Errorf("Crash while fetching the closure transactions!\nerr.Error(): %v\n", err.Error())
	}
	return nil
}

type ClosurePaymentTransactions struct {
	TransactionId string  `db:"TRANSACTION_ID"`
	LoanId        string  `db:"LOAN_ID"`
	PaymentNo     string  `db:"PAYMENT_NUMBER"`
	PMT           float64 `db:"AMOUNT_TO_PAY"` //Amount to pay
	PayedAmount   float64 `db:"TRANSACTION_AMMOUNT"`
}

type ClosurePayments struct {
	ClosureId int
	CPT       []ClosurePaymentTransactions
}

func (closurePayment *ClosurePayments) Fetch(db *sqlx.DB) error {
	query := `SELECT 
        t.TRANSACTION_ID,
        p.LOAN_ID,
        p.PAYMENT_NUMBER,
        p.AMOUNT_TO_PAY,
        t.TRANSACTION_AMMOUNT
    FROM CLOSURE_PAYMENTS cp 
    JOIN PAYMENTS p
    ON cp.PAYMENT_ID = p.PAYMENT_ID
    JOIN PAYMENT_TRANSACTIONS pt 
    ON cp.PAYMENT_ID = pt.PAYMENT_ID 
    JOIN TRANSACTIONS t 
    ON t.TRANSACTION_ID = pt.TRANSACTION_ID
    WHERE cp.CLOSURE_ID = ?`
	err := db.Select(&closurePayment.CPT, query, closurePayment.ClosureId)
	if err != nil {
		return fmt.Errorf("Crash while fetching the closure payments\nerr.Error() %v\n", err.Error())
	}
	return nil
}

type Payouts struct {
	PayoutId              int     `db:"PAYOUT_ID"`
	ClosureId             int     `db:"CLOSURE_ID"`
	AccountId             string  `db:"ACCOUNT_ID"`
	AccountBalance        float64 `db:"ACCOUNT_BALANCE"`
	ApportationPercentage int     `db:"APPORTATION_PERCENTAGE"`
	AccountProfit         float64 `db:"ACCOUNT_PROFIT"`
	DecimalPercentage     float32 `db:"DECIMAL_PERCENTAGE"`
	PayoutDate            string  `db:"PAYOUT_DATE"`
	Name                  string  `db:"NAME"`
}

func (payout *Payouts) Fetch(db *sqlx.DB) error {
	query := `SELECT *, APPORTATION_PERCENTAGE/100.00 AS decimal_percentage
        FROM payouts
        WHERE payout_id = ?`
	err := db.Get(payout, query, payout.PayoutId)
	if err != nil {
		return fmt.Errorf("Crash while fetching a dividend!\nerr.Error() %v\n", err.Error())
	}
	return nil
}

type AffiliateReports struct {
	UserId             string    `db:"USER_ID"`
	Name               string    `db:"NAME"`
	HiringDate         time.Time `db:"HIRING_DATE"`
	SavingsBalance     float64   `db:"SAVINGS_BALANCE"`
	ApportationBalance float64   `db:"APPORTATION_BALANCE"`
	Total              float64   `db:"TOTAL"`
	HiringDateFmt      string
}

func (affiliateReport *AffiliateReports) Fetch(db *sqlx.DB) error {
	query := `SELECT 
            MAX(u.USER_ID) AS user_id,
            MAX(u.FIRST_NAME) || ' ' || MAX(u.FIRST_LASTNAME) AS name,
            MAX(u.HIRING_DATE) AS hiring_date,
            MAX(CASE WHEN a.ACCOUNT_TYPE = 'CAR' THEN a.BALANCE END) AS savings_balance,
            MAX(CASE WHEN a.ACCOUNT_TYPE = 'CAP' THEN a.BALANCE END) AS apportation_balance,
            SUM(a.BALANCE) AS total
        FROM USERS u
        JOIN ACCOUNTS a
        ON u.USER_ID = a.USER_ID
        WHERE a.USER_ID = ?
        GROUP BY u.USER_ID;`
	err := db.Get(affiliateReport, query, affiliateReport.UserId)
	if err != nil {
		return fmt.Errorf("Crash while fetching account reports!\nerr.Error(): %v\n", err.Error())
	}
	affiliateReport.HiringDateFmt = affiliateReport.HiringDate.Format("2006-03-02")
	return nil
}

type LoanTransactions struct {
	LoanId        string    `db:"LOAN_ID"`
	TransactionId string    `db:"TRANSACTION_ID"`
	Amount        float64   `db:"TRANSACTION_AMMOUNT"`
	Date          time.Time `db:"TRANSACTION_DATE"`
	PaymentNo     string    `db:"PAYMENT_NUMBER"`
	FmtDate       string
}

type Accounts struct {
	UserID               string    `db:"USER_ID"`
	AccountType          string    `db:"ACCOUNT_TYPE"`
	AccountId            string    `db:"ACCOUNT_ID"`
	Balance              float64   `db:"BALANCE"`
	CreatedBy            string    `db:"CREATED_BY"`
	CreationDate         time.Time `db:"CREATION_DATE"`
	Modified             string    `db:"MODIFIED_BY"`
	LastModificationDate time.Time `db:"LAST_MODIFICATION_DATE"`
}

func (account *Accounts) Fetch(db *sqlx.DB) error {
	query := `SELECT *
        FROM ACCOUNTS c
        WHERE c.USER_ID = ?
        AND c.ACCOUNT_TYPE = ?`
	err := db.Get(account, query, account.UserID, account.AccountType)
	if err != nil {
		return fmt.Errorf("Error while fetching the account!\nerr.Error(): %v\n", err.Error())
	}
	return nil
}

type Liquidations struct {
	Id         int            `db:"LIQUIDATION_ID"`
	AccountId  string         `db:"ACCOUNT_ID"`
	Type       string         `db:"LIQUIDATION_TYPE"`
	Date       time.Time      `db:"RETIREMENT_DATE"`
	TotalMoney float64        `db:"TOTAL_MONEY"`
	Comment    sql.NullString `db:"TRANSACTION_COMMENT"`
	DateFmt    string
}

func (liquidation *Liquidations) Insert(db *sqlx.DB) error {
	if liquidation.Type == "P" {
		query := `CALL sp_create_partial_liquidation(?,?,?,?,?)`
		_, err := db.Exec(query,
			sql.Out{Dest: &liquidation.Id},
			liquidation.AccountId,
			liquidation.TotalMoney,
			liquidation.Date,
			liquidation.Comment)
		if err != nil {
			return fmt.Errorf("Crash while inserting the partial liquidation!\nerr.Error(): %v\n", err.Error())
		}
		return nil
	} else if liquidation.Type == "T" {
		query := `CALL sp_create_total_liquidation(?,?,?,?)`
		_, err := db.Exec(query,
			sql.Out{Dest: &liquidation.Id},
			liquidation.AccountId,
			liquidation.Date,
			liquidation.Comment)
		if err != nil {
			return fmt.Errorf("Crash while inserting the total liquidation!\nerr.Error(): %v\n", err.Error())
		}
		return nil
	}
	return fmt.Errorf("The liquidation type should be (T)otal or (P)artial")
}

func (liquidation *Liquidations) Update(db *sqlx.DB) error {
	return fmt.Errorf("Liquidation updates are not suppoerted at this time!")
}
