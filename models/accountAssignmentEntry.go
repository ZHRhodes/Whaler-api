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
	var entry = &AccountAssignmentEntry{
		AccountID:  newEntry.AccountID,
		AssignedBy: newEntry.AssignedBy,
		AssignedTo: newEntry.AssignedTo,
	}

	var err = db.Create(entry).Error

	if len(entry.ID) == 0 {
		return nil, err
	}
	
	db.First(&Account{}, newEntry.AccountID).Association("AssignmentEntries").Append(entry)

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
