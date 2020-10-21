package models

import (
	"fmt"

	"github.com/heroku/whaler-api/graph/model"
)

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
