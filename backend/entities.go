package backend

import (
	"database/sql"
	"fmt"
	"strings"
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

type User struct {
	UserNumber           int            `db:"USER_NUMBER"`
	UserId               string         `db:"USER_ID"`
	Password             string         `db:"PASSWORD"`
	FirstName            sql.NullString `db:"FIRST_NAME"`
	SecondName           sql.NullString `db:"SECOND_NAME"`
	FirstLastname        sql.NullString `db:"FIRST_LASTNAME"`
	SecondLastname       sql.NullString `db:"SECOND_LASTNAME"`
	AddressHouseNumber   sql.NullString `db:"ADDRESS_HOUSE_NUMBER"`
	AddressStreet        sql.NullString `db:"ADDRESS_STREET"`
	AddressAvenue        sql.NullString `db:"ADDRESS_AVENUE"`
	AddressCity          sql.NullString `db:"ADDRESS_CITY"`
	AddressDepartment    sql.NullString `db:"ADDRESS_DEPARTMENT"`
	AddressReference     sql.NullString `db:"ADDRESS_REFERENCE"`
	PrimaryEmail         sql.NullString `db:"PRIMARY_EMAIL"`
	SecondaryEmail       sql.NullString `db:"SECONDARY_EMAIL"`
	BirthDate            sql.NullTime   `db:"BIRTH_DATE"`
	HiringDate           time.Time      `db:"HIRING_DATE"`
	CreatedBy            sql.NullString `db:"CREATED_BY"`
	CreationDate         time.Time      `db:"CREATION_DATE"`
	ModifiedBy           sql.NullString `db:"MODIFIED_BY"`
	LastModificationDate time.Time      `db:"LAST_MODIFICATION_DATE"`
}

func (user *User) Insert(db *sqlx.DB) error {
	query :=
		strings.ToUpper(`SELECT user_number, user_id FROM FINAL TABLE (
      INSERT INTO users (
        password, first_name, second_name, first_lastname, second_lastname, address_house_number, address_street,
        address_avenue, address_city, address_department, address_reference, primary_email,	secondary_email, birth_date,
        hiring_date, created_by,	creation_date, modified_by,	last_modification_date
      ) VALUES (
        :password, :first_name, :second_name, :first_lastname, :second_lastname, :address_house_number, :address_street,
        :address_avenue, :address_city, :address_department, :address_reference, :primary_email, :secondary_email, :birth_date,
        :hiring_date, :created_by,	:creation_date, :modified_by,	:last_modification_date
      )
    )`)

	err := db.QueryRowx(query, &user).Scan(&user.UserNumber, &user.UserId)
	if err != nil {
		return fmt.Errorf("Crash at user insert!\nerr.Error(): %v\n", err.Error())
	}

	fmt.Println("User inserted succesfully")
	return nil
}

func (user *User) Update(db *sqlx.DB) error {
	query :=
	strings.ToUpper(`
  UPDATE users SET
    password=:password, first_name=:first_name, second_name=:second_name, first_lastname=:first_lastname, seccond_lastname=:second_lastname,
    address_house_number=:address_house_number, address_street=:address_street, address_avenue=:address_avenue, address_city=:address_city,
    address_department=:address_department, address_reference=:address_reference, primary_email=:primary_email, secondary_email=:secondary_email,
    modified_by=:modified_by, last_modification_date=:last_modification_date
  WHERE user_id=:user_id
  `)
  
	_, err := db.NamedExec(query, &user)
	if err != nil {
		return fmt.Errorf("Crash while updating user\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

func (user *User) Fetch(db *sqlx.DB) error {
	err := db.Get(user, `SELECT * FROM users WHERE user_id=? AND password=?`, user.UserId, user.Password)
	if err != nil {
		return fmt.Errorf("Crash while fetching user\nerr.Error(): %v\n", err.Error())
	}

	return nil
}

type PhoneNumbers struct {
	user_id           string
	user_phone_number int
	region_number     int
}

