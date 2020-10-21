package models

import (
	"fmt"

	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
)

type Contact struct {
	DBModel
	FirstName             string                     `json:"firstName"`
	LastName              string                     `json:"lastName"`
	JobTitle              string                     `json:"jobTitle"`
	State                 string                     `json:"state"`
	Email                 string                     `json:"email"`
	Phone                 string                     `json:"phone"`
	AccountID			  string					 `json:"accountID"` //TODO.. should i do this, or add an Account prop to Contact? need to connect them
	AssignmentEntries     []ContactAssignmentEntry   `json:"assignmentEntries" gorm:"foreignKey:ContactID;references:ID"`
	// Account               Account                 `json:"account"`
	// Seniority             string                     `json:"seniority"`
	// Persona               string                     `json:"persona"`
	// AssignedTo            User                    `json:"assignedTo"`
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
	contact := &Contact{
		FirstName: newContact.FirstName,
		LastName:  newContact.LastName,
		JobTitle:  *newContact.JobTitle,
		State:     *newContact.State,
		Email:     *newContact.Email,
		Phone:     *newContact.Phone,
		AccountID: *newContact.AccountID,
		//figure out how to use AccountID to tie this contact to an account in db
	}

	err := DB().Create(contact).Error

	if len(contact.ID) == 0 {
		return nil, err
	}

	return contact, nil
}

type ContactAssignmentEntry struct {
	DBModel
	ContactID  string  `json:"contactId"`
	AssignedBy string  `json:"assignedBy"`
	AssignedTo *string `json:"assignedTo"`
}

func CreateContactAssignmentEntry(newEntry model.NewContactAssignmentEntry) (*ContactAssignmentEntry, error) {
	entry := &ContactAssignmentEntry{
		ContactID:  newEntry.ContactID,
		AssignedBy: newEntry.AssignedBy,
		AssignedTo: newEntry.AssignedTo,
	}
	
	// err = db.Debug().Model(&Contact{}).Where("id = ?", newEntry.ContactID).Update("LatestAssignmentEntry", entry).Error
	// err := db.Model(&Contact{}).Where("id = ?", newEntry.ContactID).Association("Languages").Order("createdDate desc").Find(&languages)

	err := db.First(&Contact{}, newEntry.ContactID).Association("AssignmentEntries").Append(entry)

	if err != nil {
		fmt.Println(err)
	}

	if len(entry.ID) == 0 {
		return nil, err
	}

	return entry, nil
}

func FetchContactAssignmentEntries(contactID string) ([]*ContactAssignmentEntry, error) {
	entries := []*ContactAssignmentEntry{}

	err := DB().Where("contact_id <> ?", contactID).Find(&entries).Error

	if err != nil {
		return []*ContactAssignmentEntry{}, err
	}

	return entries, err
}
