package backend

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ibmdb/go_ibm_db"
)

const CONNECTION_STRING string = "HOSTNAME=localhost;DATABASE=coopdb;PORT=51000;UID=db2inst1;PWD=coop4312";

func TestConnection() {
	db, err := sql.Open("go_ibm_db", CONNECTION_STRING)
  defer db.Close()

  if err != nil {
		log.Fatalln(err)
	  return ;
  }

  fmt.Println("connection succeeded")
}

