package models

import (
	"fmt"

	"github.com/heroku/whaler-api/graph/model"
)

type AccountAssignmentEntry struct {
	DBModel
	AccountID  string  `json:"accountId"`
	AssignedBy string  `json:"assignedBy"`
	AssignedTo *string `json:"assignedTo"`
}

func CreateAccountAssignmentEntry(newEntry model.NewAccountAssignmentEntry) (*AccountAssignmentEntry, error) {
	fmt.Printf("\nCreating account assignment entry")
	var entry = &AccountAssignmentEntry{
		AccountID:  newEntry.AccountID,
		AssignedBy: newEntry.AssignedBy,
		AssignedTo: newEntry.AssignedTo,
	}

	var err = db.Create(entry).Error

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var account Account
	db.First(&account, "id = ?", newEntry.AccountID).Association("AssignmentEntries").Append(entry)
	db.Model(&account).Update("AssignedTo", entry.AssignedTo)
	fmt.Printf("\nUpdating assigned to field for accountId %s to assignedTo %s", account.ID, *entry.AssignedTo)

	return entry, nil
}

func FetchAccountAssignmentEntries(accountID string) ([]*AccountAssignmentEntry, error) {
	entries := []*AccountAssignmentEntry{}

	var account Account
	var err = DB().Debug().First(&account, "id = ?", accountID).Error
	association := DB().Model(&account).Association("AssignmentEntries")
	association.Find(&entries)

	if err != nil {
		fmt.Println("Something bad happened here...")
		fmt.Println(err)
	}

	return entries, nil
}
