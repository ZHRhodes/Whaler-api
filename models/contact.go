package models

import "github.com/heroku/whaler-api/utils"

type Contact struct {
	DBModel
	firstName  string
	lastName   string
	state      string
	account    Account
	jobTitle   string
	seniority  string
	persona    string
	email      string
	phone      []string
	assignedTo User
	//notes
}

func (contact *Contact) Create() map[string]interface{} {
	DB().Create(contact)

	if contact.ID <= 0 {
		return utils.Message(5001, "Failed to create contact, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"contact": contact}
	response := utils.Message(2000, "Contact has been created", false, data)
	return response
}
