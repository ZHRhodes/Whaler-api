package models

import (
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
	var err = db.Create(entry).Error

	if len(entry.ID) == 0 {
		return nil, err
	}

	db.First(&Contact{}, newEntry.ContactID).Association("AssignmentEntries").Append(entry)

	return entry, nil
}

func FetchContactAssignmentEntries(contactID string) ([]*ContactAssignmentEntry, error) {
	entries := []*ContactAssignmentEntry{}

	DB().First(&Contact{}, contactID).Association("AssignmentEntries").Find(&entries)

	return entries, nil
}
