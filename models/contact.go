package models

import (
	"fmt"

	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
	"gorm.io/gorm/clause"
)

type Contact struct {
	DBModel
	FirstName         string                   `json:"firstName"`
	LastName          string                   `json:"lastName"`
	SalesforceID      *string                  `json:"salesforceID"`
	JobTitle          *string                  `json:"jobTitle"`
	State             *string                  `json:"state"`
	Email             *string                  `json:"email"`
	Phone             *string                  `json:"phone"`
	AccountID         *string                  `json:"accountID"`
	AssignmentEntries []ContactAssignmentEntry `json:"assignmentEntries" gorm:"foreignKey:ContactID;references:ID"`
	AssignedTo        *string                  `json:"assignedTo"`
	// Account               Account                 `json:"account"`
	// Seniority             string                     `json:"seniority"`
	// Persona               string                     `json:"persona"`
	// ExternalID            string                     `json:"externalID"`
	//notes
}

func (contact *Contact) Create() map[string]interface{} {
	DB().Create(contact)

	if len(contact.ID) == 0 {
		return utils.Message(5001, "Failed to create contact, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"contact": contact}
	response := utils.Message(2000, "Contact has been created", false, data)
	return response
}

func CreateContact(newContact model.NewContact) (*Contact, error) {
	contact := createContactFromNewContact(newContact)

	err := DB().Create(contact).Error

	if len(contact.ID) == 0 {
		return nil, err
	}

	return contact, nil
}

func SaveContacts(newContacts []*model.NewContact) ([]*Contact, error) {
	var contacts = []*Contact{}
	for _, newContact := range newContacts {
		contacts = append(contacts, createContactFromNewContact(*newContact))
	}

	err := DB().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "first_name", "last_name", "job_title",
			"salesforce_id", "state", "email", "phone", "account_id", "assigned_to"}),
	}).Create(&contacts).Error

	if err != nil {
		fmt.Println("erroring saving contacts.", err)
	}

	//This might be excessive, especially if FE can't perfectly hide the reloading
	if len(newContacts) > 0 && newContacts[0] != nil {
		go Consumer.ModelChanged(*newContacts[0].AccountID)
	}

	return contacts, err
}

func FetchContacts(accountID string) ([]*Contact, error) {
	var contacts = []*Contact{}
	err := db.Model(&Account{DBModel: DBModel{ID: accountID}}).Association("Contacts").Find(&contacts)
	return contacts, err
}

func createContactFromNewContact(newContact model.NewContact) *Contact {
	var id string
	if newContact.ID != nil {
		id = *newContact.ID
	}
	return &Contact{
		DBModel:      DBModel{ID: id},
		FirstName:    newContact.FirstName,
		LastName:     newContact.LastName,
		SalesforceID: newContact.SalesforceID,
		JobTitle:     newContact.JobTitle,
		State:        newContact.State,
		Email:        newContact.Email,
		Phone:        newContact.Phone,
		AccountID:    newContact.AccountID,
		AssignedTo:   newContact.AssignedTo,
	}
}
