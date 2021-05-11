package models

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type Note struct {
	DBModel
	OwnerID   string `json:"-" gorm:"index:owner_account"`
	AccountID string `json:"accountId" gorm:"index:owner_account"`
	Content   string `json:"content"`
}

func FetchNote(ownerID string, accountID string) (*Note, error) {
	// query := map[string]interface{}{"owner_id": ownerID, "account_id": accountID}
	query := map[string]interface{}{"account_id": accountID}
	var note Note
	err := db.Table("notes").Where(query).FirstOrInit(&note, Note{OwnerID: ownerID, AccountID: accountID, Content: ""}).Error
	return &note, err
}

func SaveNote(ownerID string, note Note) (*Note, error) {
	fmt.Printf("\nOwnerId %s is saving note with id %s and accountId %s\n", ownerID, note.ID, note.AccountID)
	note.OwnerID = ownerID
	conflictColumn := []clause.Column{{Name: "id"}}
	assignmentColumn := clause.AssignmentColumns([]string{"updated_at", "account_id", "owner_id", "content"})
	err := db.Clauses(clause.OnConflict{Columns: conflictColumn, DoUpdates: assignmentColumn}).Create(&note).Error
	return &note, err
}

func SaveNoteContent(accountID string, newContent string) error {
	err := DB().Model(&Note{}).Where("account_id = ?", accountID).Update("Content", newContent).Error
	if err != nil {
		fmt.Printf("\nFailed to save new content to note with accountId %s", accountID)
	}
	return err
}
