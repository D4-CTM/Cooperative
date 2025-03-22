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
	query := `CALL sp_insert_user(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
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
		user.LastModificationDate)

	if err != nil {
		return fmt.Errorf("Crash at user insert!\nerr.Error(): %v\n", err.Error())
	}

	fmt.Println("User inserted succesfully, id:", user.UserId)
	return nil
}

func (user *Users) Update(db *sqlx.DB) error {
	query := `CALL sp_update_user(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
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
		user.LastModificationDate)

	if err != nil {
		return fmt.Errorf("Crash while updating user\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

func (user *Users) Fetch(db *sqlx.DB) error {
	err := db.Get(user, `SELECT * FROM users WHERE user_id=? AND password=?`, user.UserId, user.Password)
	if err != nil {
		return fmt.Errorf("Crash while fetching user\nerr.Error(): %v\n", err.Error())
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
		return err
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
	LoanId   string    `db:"LOAN_ID"`
	UserId   string    `db:"USER_ID"`
	Periods  int       `db:"LOAN_PERIODS"`
	Interest float32   `db:"LOAN_INTEREST"`
	Capital  float64   `db:"REQUESTED_AMOUNT"`
	Date     time.Time `db:"LOAN_DATE"`
	IsPayed  bool      `db:"IS_PAYED"`
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
            PT.Payment.PaymentId)

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

