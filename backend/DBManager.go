package backend

import (
	"fmt"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/jmoiron/sqlx"
)

const CONNECTION_STRING string = "HOSTNAME=localhost;DATABASE=coopdb;PORT=51000;UID=db2inst1;PWD=coop4312"

func getConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("go_ibm_db", CONNECTION_STRING)
	if err != nil {
		return nil, fmt.Errorf("Crashed on connection!\nerr.Error():%v\n", err.Error())
	}
	return db, nil
}

func Fetch(fetchable Fetchable) error {
	con, err := getConnection()
	defer con.Close()

	if err != nil {
		return err
	}

	return fetchable.Fetch(con)
}
