package models

import "github.com/heroku/whaler-api/utils"

type Contact struct {
	DBModel
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	State      string  `json:"state"`
	Account    Account `json:"account"`
	JobTitle   string  `json:"jobTitle"`
	Seniority  string  `json:"seniority"`
	Persona    string  `json:"persona"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	AssignedTo User    `json:"assignedTo"`
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

func CreateContact(newContact model.NewContact) (*Contact, error) {
	contact := Contact{
		FirstName: newContact.FirstName,
		LastName: newContact.LastName,
		State: newContact.State,
		//figure out how to use AccountID to tie this contact to an account in db
		JobTitle: newContact.JobTitle,
		Seniority: newContact.Seniority,
		Peronsa: newContact.Persona,
		Email: newContact.Email,
		Phone: newContact.Phone
	}

	err := DB().Create(contact).Error

	if account.ID <= 0 {
		return nil, err
	}

	return contact, nil
}