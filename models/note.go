package models

import (
	"gorm.io/gorm/clause"
)

type Note struct {
	DBModel
	OwnerID   string `json:"-" gorm:"index:owner_account"`
	AccountID string `json:"accountID" gorm:"index:owner_account"`
	Content   string `json:"content"`
}

func FetchNote(ownerID string, accountID string) (*Note, error) {
	query := map[string]interface{}{"owner_id": ownerID, "account_id": accountID}
	var note Note
	err := db.Table("Notes").Where(query).First(&note).Error
	return &note, err
}

func SaveNote(note Note) (*Note, error) {
	conflictColumn := []clause.Column{{Name: "id"}}
	assignmentColumn := clause.AssignmentColumns([]string{"updated_at", "account_id", "owner_id", "content"})
	err := db.Clauses(clause.OnConflict{Columns: conflictColumn, DoUpdates: assignmentColumn}).Create(&note).Error
	return &note, err
}
