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
	entry := &AccountAssignmentEntry{
		AccountID:  newEntry.AccountID,
		AssignedBy: newEntry.AssignedBy,
		AssignedTo: newEntry.AssignedTo,
	}

	err := db.First(&Account{}, newEntry.AccountID).Association("AssignmentEntries").Append(entry)

	if err != nil {
		fmt.Println(err)
	}

	if len(entry.ID) == 0 {
		return nil, err
	}

	return entry, nil
}

func FetchAccountAssignmentEntries(accountID string) ([]*AccountAssignmentEntry, error) {
	entries := []*AccountAssignmentEntry{}

	err := DB().Where("account_id <> ?", accountID).Find(&entries).Error

	if err != nil {
		return []*AccountAssignmentEntry{}, err
	}

	return entries, err
}
