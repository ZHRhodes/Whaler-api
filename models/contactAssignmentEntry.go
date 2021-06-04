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

func CreateContactAssignmentEntry(senderID *string, newEntry model.NewContactAssignmentEntry) (*ContactAssignmentEntry, error) {
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

	var contact Contact

	err = db.First(&contact, "id = ?", newEntry.ContactID).Association("AssignmentEntries").Append(entry)

	if err != nil {
		fmt.Printf("\nError adding assignment entry for contact %s", contact.ID)
	}

	err = db.Model(&contact).Update("AssignedTo", entry.AssignedTo).Error

	if err != nil {
		fmt.Printf("\nError updating assignedTo on contact with id %s", contact.ID)
	}

	fmt.Printf("\nUpdating assigned to field for contactId %s to assignedTo %s", contact.ID, *entry.AssignedTo)

	if contact.AccountID != nil {
		go Consumer.ModelChanged(*contact.AccountID, senderID)
	}

	return entry, nil
}

func FetchContactAssignmentEntries(contactID string) ([]*ContactAssignmentEntry, error) {
	entries := []*ContactAssignmentEntry{}

	var contact Contact
	var err = DB().Debug().First(&contact, "id = ?", contactID).Error
	association := DB().Model(&contact).Association("AssignmentEntries")
	association.Find(&entries)

	if err != nil {
		fmt.Println("Something bad happened here...")
		fmt.Println(err)
	}

	return entries, nil
}
