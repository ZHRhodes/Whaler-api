package models

import "github.com/heroku/whaler-api/utils"

type Contact struct {
	DBModel
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	State      string   `json:"state"`
	Account    Account  `json:"account"`
	JobTitle   string   `json:"jobTitle"`
	Seniority  string   `json:"seniority"`
	Persona    string   `json:"persona"`
	Email      string   `json:"email"`
	Phone      []string `json:"phone" gorm:"type:varchar(64)[]"`
	AssignedTo User     `json:"assignedTo"`
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
