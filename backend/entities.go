package backend

import "time"

type User struct {
	userNumber int32
	userId string
  password string
  firstName string
  secondName string
  firstLastname string
  secondLastname string
	addressHouse_number string
	addressStreet string
	addressAvenue string
	addressCity string
	addressDepartment string
	addressReference string
	primaryEmail string
	secondaryEmail string
	birthDate time.Time
	hiringDate time.Time 
	createdBy string
	creationDate time.Time
	modifiedBy string
	lastModificationDate time.Time
}

type PhoneNumbers struct {
  user_id string
	user_phone_number int
	region_number int
}

