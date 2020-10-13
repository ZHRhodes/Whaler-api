package models

import (
	"fmt"

	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
)

type Contact struct {
	DBModel
	FirstName             string                 `json:"firstName"`
	LastName              string                 `json:"lastName"`
	State                 string                 `json:"state"`
	// Account               Account                `json:"account"`
	JobTitle              string                 `json:"jobTitle"`
	Seniority             string                 `json:"seniority"`
	Persona               string                 `json:"persona"`
	Email                 string                 `json:"email"`
	Phone                 string                 `json:"phone"`
	// AssignedTo            User                   `json:"assignedTo"`
	ExternalIDSon            string                 `json:"externalID"`
	LatestAssignmentEntry ContactAssignmentEntry `json:"latestAssignmentEntry" gorm:"foreignKey:ContactID;references:ID"`
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
	contact := &Contact{
		FirstName: newContact.FirstName,
		LastName:  newContact.LastName,
		State:     *newContact.State,
		//figure out how to use AccountID to tie this contact to an account in db
		JobTitle:  *newContact.JobTitle,
		Seniority: *newContact.Seniority,
		Email:     *newContact.Email,
		Phone:     *newContact.Phone,
	}

	err := DB().Create(contact).Error

	if contact.ID <= 0 {
		return nil, err
	}

	return contact, nil
}

type ContactAssignmentEntry struct {
	DBModel
	ContactID  int     `json:"contactId"`
	AssignedBy string  `json:"assignedBy"`
	AssignedTo *string `json:"assignedTo"`
}

func CreateContactAssignmentEntry(newEntry model.NewContactAssignmentEntry) (*ContactAssignmentEntry, error) {
	entry := &ContactAssignmentEntry{
		ContactID:  newEntry.ContactID,
		AssignedBy: newEntry.AssignedBy,
		AssignedTo: newEntry.AssignedTo,
	}

	var err = db.Debug().Create(entry).Error

	if err != nil {
		fmt.Println(err)
	}

	err = db.Debug().Model(&Contact{}).Where("id = ?", newEntry.ContactID).Update("latestAssignmentEntry", entry).Error

	if err != nil {
		fmt.Println(err)
	}

	if entry.ID <= 0 {
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
