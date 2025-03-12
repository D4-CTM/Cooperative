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

type User struct {
	UserId               string         `db:"USER_ID"`
	Password             string         `db:"PASSWORD"`
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

func (user *User) Insert(db *sqlx.DB) error {
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

func (user *User) Update(db *sqlx.DB) error {
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
